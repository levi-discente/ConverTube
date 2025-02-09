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
import { diskStorage } from 'multer';
import * as path from 'path';
import * as fs from 'fs';
import { Response } from 'express';

@Controller('conversion')
export class ConversionController {
  constructor(private readonly conversionService: ConversionService) {}

  @Post('upload')
  @UseInterceptors(
    FileInterceptor('file', {
      storage: diskStorage({
        destination: './uploads',
        filename: (req, file, callback) => {
          const uniqueSuffix =
            Date.now() + '-' + Math.round(Math.random() * 1e9);
          const ext = path.extname(file.originalname);
          const originalName = path.basename(file.originalname, ext);
          const filename = `${originalName}-${uniqueSuffix}${ext}`;
          callback(null, filename);
        },
      }),
    }),
  )
  async uploadFile(
    @UploadedFile() file: Express.Multer.File,
    @Body('outputFormat') outputFormat: string,
    @Body('quality') quality: string,
  ): Promise<{ operationId: string; responseQueue: string }> {
    const filenameWithExtension = file.filename;
    const { name } = path.parse(filenameWithExtension);
    const filePath = await this.conversionService.storeFile(file);
    const { operationId, responseQueue } =
      await this.conversionService.createConversionJob(
        filePath,
        name,
        outputFormat,
        quality,
      );
    return { operationId, responseQueue };
  }

  @Get('download/:filename')
  downloadFile(@Param('filename') filename: string, @Res() res: Response) {
    const filePath = path.join(__dirname, '..', '..', 'uploads', filename);

    if (!fs.existsSync(filePath)) {
      throw new NotFoundException('Arquivo não encontrado');
    }

    res.download(filePath, filename);
  }
}
