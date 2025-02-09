/* eslint-disable @typescript-eslint/no-floating-promises */
import { Injectable, InternalServerErrorException } from '@nestjs/common';
import { Client } from 'minio';
import { Readable } from 'stream';

@Injectable()
export class StorageService {
  private minioClient: Client;
  private bucketName = process.env.MINIO_BUCKET || 'uploads';

  constructor() {
    this.minioClient = new Client({
      endPoint:
        process.env.MINIO_ENDPOINT ||
        'minio-headless.default.svc.cluster.local',
      port: Number(process.env.MINIO_PORT) || 9000,
      useSSL: process.env.MINIO_USE_SSL === 'true',
      accessKey: process.env.MINIO_ACCESS_KEY || 'minio',
      secretKey: process.env.MINIO_SECRET_KEY || 'minio123',
    });

    void this.ensureBucketExists();
  }

  async ensureBucketExists(): Promise<void> {
    try {
      const exists = await this.minioClient.bucketExists(this.bucketName);
      if (!exists) {
        await this.minioClient.makeBucket(this.bucketName, 'pe-east-1');
        console.log(`Bucket ${this.bucketName} criado com sucesso`);
      }
    } catch (err) {
      console.error('Erro ao verificar/criar bucket:', err);
    }
  }

  async uploadFile(
    fileBuffer: Buffer,
    objectName: string,
    fileSize: number,
  ): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      this.minioClient.putObject(
        this.bucketName,
        objectName,
        fileBuffer,
        fileSize,
        (err, etag) => {
          if (err) {
            return reject(
              new InternalServerErrorException(
                'Erro ao enviar arquivo para o MinIO',
              ),
            );
          }
          console.log('ETag:', etag);
          resolve(objectName);
        },
      );
    });
  }

  async getFileStream(objectName: string): Promise<Readable> {
    try {
      return await this.minioClient.getObject(this.bucketName, objectName);
    } catch (err) {
      throw new InternalServerErrorException(
        'Erro ao recuperar o arquivo do MinIO: ' + err,
      );
    }
  }
}
