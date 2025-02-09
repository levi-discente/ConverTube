import {
  WebSocketGateway,
  WebSocketServer,
  OnGatewayInit,
  OnGatewayConnection,
  OnGatewayDisconnect,
} from "@nestjs/websockets";
import { Server, Socket } from "socket.io";
import * as amqp from "amqplib";
import { Logger } from "@nestjs/common";
import { ResponseMessage } from "./conversion.interfaces";

@WebSocketGateway({
  path: "/ws/",
  cors: {
    origin: "*",
    methods: ["GET", "POST"],
    allowedHeaders: ["Content-Type", "Authorization"],
    credentials: true,
  },
})
export class ConversionGateway
  implements OnGatewayInit, OnGatewayConnection, OnGatewayDisconnect
{
  @WebSocketServer()
  server: Server;

  private readonly logger = new Logger(ConversionGateway.name);
  private connection: amqp.Connection;
  private channel: amqp.Channel;
  private activeConsumers: Map<string, string> = new Map(); // Mapeia operationId → consumerTag

  async afterInit(server: Server) {
    try {
      this.connection = await amqp.connect(
        process.env.RABBITMQ_URL || "amqp://guest:guest@localhost:5672/",
      );
      this.channel = await this.connection.createChannel();
      this.logger.log("RabbitMQ connection established in WebSocket Gateway");
    } catch (err) {
      this.logger.error("Failed to connect to RabbitMQ", err);
    }
  }

  handleConnection(client: Socket) {
    const operationId = client.handshake.query.operationId as string;

    if (!operationId) {
      this.logger.warn(`Client ${client.id} connected without operationId`);
      client.disconnect();
      return;
    }

    this.logger.log(
      `Client ${client.id} connected for operationId: ${operationId}`,
    );
    this.subscribeToJob(client, operationId);
  }

  handleDisconnect(client: Socket) {
    this.logger.log(`Client disconnected: ${client.id}`);
  }

  async subscribeToJob(client: Socket, operationId: string) {
    const queueName = `conversion_responses_${operationId}`;

    try {
      await this.channel.assertQueue(queueName, {
        durable: true,
        autoDelete: true,
        exclusive: false,
      });

      const { consumerTag } = await this.channel.consume(
        queueName,
        (msg) => {
          if (msg) {
            try {
              const response: ResponseMessage = JSON.parse(
                msg.content.toString(),
              );
              client.emit("jobUpdate", response);
              this.logger.log(
                `Sent update to client ${client.id}: ${JSON.stringify(response)}`,
              );

              // Finaliza se for erro ou sucesso
              if (
                response.status === "error" ||
                response.status === "success"
              ) {
                this.logger.log(`Finalizando job ${operationId}`);
                this.cleanupJob(operationId);
                client.disconnect();
              }
            } catch (error) {
              this.logger.error("Failed to parse message", error);
            }
          }
        },
        { noAck: true },
      );

      // Salva referência ao consumidor para futura remoção
      this.activeConsumers.set(operationId, consumerTag);
      this.logger.log(`Subscribed to job responses on queue: ${queueName}`);
    } catch (err) {
      this.logger.error(`Failed to subscribe to queue ${queueName}`, err);
      client.emit("jobUpdate", {
        status: "error",
        message: "Subscription failed",
      });
    }
  }

  cleanupJob(operationId: string) {
    const consumerTag = this.activeConsumers.get(operationId);
    if (consumerTag) {
      this.channel
        .cancel(consumerTag)
        .catch((err) => this.logger.error("Error canceling consumer", err));
      this.activeConsumers.delete(operationId);
    }
  }
}
