import {
  Injectable,
  Logger,
  OnModuleDestroy,
  OnModuleInit,
} from '@nestjs/common';
import * as amqp from 'amqplib';

@Injectable()
export class RabbitMQPublisher implements OnModuleInit, OnModuleDestroy {
  private connection: amqp.Connection;
  private channel: amqp.Channel;
  private readonly logger = new Logger(RabbitMQPublisher.name);
  private readonly rabbitMQUrl =
    process.env.RABBITMQ_URL ||
    'amqp://guest:guest@rabbitmq.default.svc.cluster.local:5672/';

  async onModuleInit() {
    this.connection = await amqp.connect(this.rabbitMQUrl);
    this.channel = await this.connection.createChannel();
    this.logger.log('RabbitMQ connection and channel established');
  }

  async onModuleDestroy() {
    await this.channel.close();
    await this.connection.close();
  }

  async publish(queue: string, message: any): Promise<void> {
    await this.channel.assertQueue(queue, {
      durable: true,
      autoDelete: true,
    });

    const buffer = Buffer.from(JSON.stringify(message));
    const ok = this.channel.sendToQueue(queue, buffer, { persistent: true });
    if (ok) {
      this.logger.log(`Message published to queue ${queue}`);
    } else {
      this.logger.error(`Failed to publish message to queue ${queue}`);
    }
  }
}
