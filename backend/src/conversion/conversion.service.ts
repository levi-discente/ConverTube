import { Injectable, Logger, BadRequestException } from '@nestjs/common';
import { v4 as uuidv4 } from 'uuid';
import { RabbitMQPublisher } from '../rabbitmq/rabbitmq.publisher';
import * as path from 'path';

@Injectable()
export class ConversionService {
  private readonly logger = new Logger(ConversionService.name);

  constructor(private readonly rabbitPublisher: RabbitMQPublisher) {}

  private readonly supportedFormats: Record<string, boolean> = {
    mp4: true,
    avi: true,
    mkv: true,
    mov: true,
    flv: true,
    webm: true,
    ogg: true,
    wav: true,
    mp3: true,
    aac: true,
    flac: true,
    wma: true,
    gif: true,
  };

  private readonly videoFormats = new Set([
    'mp4',
    'avi',
    'mkv',
    'mov',
    'flv',
    'webm',
  ]);
  private readonly audioFormats = new Set([
    'ogg',
    'wav',
    'mp3',
    'aac',
    'flac',
    'wma',
  ]);
  private readonly imageFormats = new Set(['gif']);

  async storeFile(file: Express.Multer.File): Promise<string> {
    const ext = this.getFileExtension(file.originalname);

    if (!this.supportedFormats[ext]) {
      throw new BadRequestException(
        `Formato de arquivo '${ext}' não é suportado.`,
      );
    }

    return Promise.resolve(path.resolve(file.path));
  }

  async createConversionJob(
    filePath: string,
    fileName: string,
    outputFormat: string,
    quality: string,
  ): Promise<{ operationId: string; responseQueue: string }> {
    const inputFormat = this.getFileExtension(filePath);

    if (!this.supportedFormats[outputFormat]) {
      throw new BadRequestException(
        `Formato de saída '${outputFormat}' não é suportado.`,
      );
    }

    if (this.isInvalidConversion(inputFormat, outputFormat)) {
      throw new BadRequestException(
        `Não é possível converter ${inputFormat.toUpperCase()} para ${outputFormat.toUpperCase()}.`,
      );
    }

    const operationId = uuidv4();
    const responseQueue = `conversion_responses_${operationId}`;
    const job = {
      operation_id: operationId,
      file_path: filePath,
      file_name: fileName,
      output_format: outputFormat,
      request_time: new Date().toISOString(),
      quality: quality,
      response_queue: responseQueue,
    };

    await this.rabbitPublisher.publish(`conversion_jobs_${operationId}`, job);
    this.logger.log(`Conversion job created: ${operationId}`);

    return { operationId, responseQueue };
  }

  private getFileExtension(filename: string): string {
    return filename.split('.').pop()?.toLowerCase() || '';
  }

  private isInvalidConversion(
    inputFormat: string,
    outputFormat: string,
  ): boolean {
    if (
      this.audioFormats.has(inputFormat) &&
      this.videoFormats.has(outputFormat)
    ) {
      return true; // Áudio → Vídeo não permitido
    }
    if (
      this.imageFormats.has(inputFormat) &&
      this.videoFormats.has(outputFormat)
    ) {
      return true; // GIF → Vídeo não permitido
    }
    return false; // Caso contrário, conversão válida
  }
}
