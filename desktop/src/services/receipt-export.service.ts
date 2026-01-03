import { Injectable } from "@angular/core";
import { take, tap } from "rxjs";
import { ExportFormat, ExportService, ReceiptPagedRequestCommand } from "../open-api/index";
import { downloadFile } from "../utils/file";

@Injectable({
  providedIn: "root"
})
export class ReceiptExportService {

  constructor(
    private exportService: ExportService,
  ) { }

  public exportReceiptsFromFilter(groupId: string, filter: ReceiptPagedRequestCommand): void {
    const groupIdInt = Number.parseInt(groupId);
    this.exportService.exportReceiptsForGroup(
      ExportFormat.Csv,
      groupIdInt, { ...filter, page: 1, pageSize: -1 }
    )
      .pipe(
        take(1),
        tap((blob) => {
          downloadFile(blob, "data.zip");
        })
      )
      .subscribe();
  }

  public exportReceiptsById(receiptIds: number[]): void {
    this.exportService.exportReceiptsById(ExportFormat.Csv, receiptIds)
      .pipe(
        take(1),
        tap((blob) => {
          downloadFile(blob, "data.zip");
        })
      )
      .subscribe();
  }
}
