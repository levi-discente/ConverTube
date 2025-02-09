import { Module } from "@nestjs/common";
import { ConversionController } from "./conversion.controller";
import { ConversionService } from "./conversion.service";
import { RabbitMQPublisher } from "../rabbitmq/rabbitmq.publisher";

@Module({
  controllers: [ConversionController],
  providers: [ConversionService, RabbitMQPublisher],
})
export class ConversionModule {}
