import { Injectable } from "@angular/core";
import { take, tap } from "rxjs";
import { ExportFormat, ExportService, ReceiptPagedRequestCommand } from "../open-api/index";
import { RECEIPT_DATE_FILTER_KEYS } from "../constants/receipt-filter-fields.constant";
import { toDateWireFilter } from "../utils/date-wire-filter";
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
      groupIdInt,
      {
        ...filter,
        // The picked calendar day, not a browser-local instant: see toDateWireFilter.
        filter: toDateWireFilter(filter.filter, RECEIPT_DATE_FILTER_KEYS),
        page: 1,
        pageSize: -1,
      }
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
