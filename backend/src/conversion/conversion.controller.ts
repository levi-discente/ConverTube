import {
  Controller,
  Post,
  UploadedFile,
  UseInterceptors,
  Body,
  Get,
  Param,
  Res,
  NotFoundException,
} from '@nestjs/common';
import { FileInterceptor } from '@nestjs/platform-express';
import { ConversionService } from './conversion.service';
import { memoryStorage } from 'multer';
import { Response } from 'express';
import { StorageService } from 'src/storage/storage.service';
import { Logger } from '@nestjs/common';

@Controller('conversion')
export class ConversionController {
  private readonly logger = new Logger(ConversionController.name);
  constructor(
    private readonly conversionService: ConversionService,
    private readonly storageService: StorageService,
  ) {}

  @Post('upload')
  @UseInterceptors(
    FileInterceptor('file', {
      storage: memoryStorage(),
    }),
  )
  async uploadFile(
    @UploadedFile() file: Express.Multer.File,
    @Body('outputFormat') outputFormat: string,
    @Body('quality') quality: string,
  ): Promise<{ operationId: string; responseQueue: string }> {
    console.log('Body: ', file);
    const originalName = file.originalname;
    const filePath = await this.conversionService.storeFile(file);
    const { operationId, responseQueue } =
      await this.conversionService.createConversionJob(
        filePath,
        originalName,
        outputFormat,
        quality,
      );
    return { operationId, responseQueue };
  }

  @Get('download/:filename')
  async downloadFile(
    @Param('filename') filename: string,
    @Res() res: Response,
  ) {
    try {
      const fileStream = await this.storageService.getFileStream(filename);
      fileStream.pipe(res);
    } catch (error) {
      throw new NotFoundException('Arquivo não encontrado no MinIO ' + error);
    }
  }
}
