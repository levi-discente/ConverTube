import { NestFactory } from "@nestjs/core";
import { AppModule } from "./app.module";
import { Logger } from "@nestjs/common";
import { IoAdapter } from "@nestjs/platform-socket.io";

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.enableCors({
    origin: "*",
    methods: "GET,HEAD,PUT,PATCH,POST,DELETE",
    allowedHeaders: "Content-Type,Authorization",
  });
  app.useWebSocketAdapter(new IoAdapter(app));
  const port = process.env.PORT || 3000;
  await app.listen(port);
  Logger.log(`Application listening on port ${port}`);
}
bootstrap();
