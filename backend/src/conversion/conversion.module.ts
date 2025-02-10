import { Module } from '@nestjs/common';
import { ConversionController } from './conversion.controller';
import { ConversionService } from './conversion.service';
import { ConversionGateway } from './conversion.gateway';
import { RabbitMQPublisher } from '../rabbitmq/rabbitmq.publisher';
import { KubernetesModule } from '../kubernetes/kubernetes.module';
import { StorageService } from 'src/storage/storage.service';

@Module({
  imports: [KubernetesModule],
  controllers: [ConversionController],
  providers: [
    ConversionService,
    RabbitMQPublisher,
    ConversionGateway,
    KubernetesModule,
    StorageService,
  ],
})
export class ConversionModule {}
