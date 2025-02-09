import {
  Controller,
  Post,
  UploadedFile,
  UseInterceptors,
  Body,
} from "@nestjs/common";
import { FileInterceptor } from "@nestjs/platform-express";
import { ConversionService } from "./conversion.service";
import { diskStorage } from "multer";
import * as path from "path";

@Controller("conversion")
export class ConversionController {
  constructor(private readonly conversionService: ConversionService) {}

  @Post("upload")
  @UseInterceptors(
    FileInterceptor("file", {
      storage: diskStorage({
        destination: "./uploads",
        filename: (req, file, callback) => {
          const uniqueSuffix =
            Date.now() + "-" + Math.round(Math.random() * 1e9);
          const ext = path.extname(file.originalname);
          const filename = file.fieldname + "-" + uniqueSuffix + ext;
          callback(null, filename);
        },
      }),
    }),
  )
  async uploadFile(
    @UploadedFile() file: Express.Multer.File,
    @Body("outputFormat") outputFormat: string,
    @Body("quality") quality: string,
  ): Promise<{ operationId: string; responseQueue: string }> {
    const filePath = await this.conversionService.storeFile(file);
    const { operationId, responseQueue } =
      await this.conversionService.createConversionJob(
        filePath,
        outputFormat,
        quality,
      );
    return { operationId, responseQueue };
  }
}
