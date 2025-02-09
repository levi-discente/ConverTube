import { Module } from "@nestjs/common";
import { ConversionModule } from "./conversion/conversion.module";
import { ConversionGateway } from "./conversion/conversion.gateway";

@Module({
  imports: [ConversionModule],
  providers: [ConversionGateway],
})
export class AppModule {}
