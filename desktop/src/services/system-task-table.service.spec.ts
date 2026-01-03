import { TestBed } from "@angular/core/testing";

import { NgxsModule } from "@ngxs/store";
import { SystemTaskTableState } from "../store/system-task-table.state";
import { SystemTaskTableService } from "./system-task-table.service";

describe("SystemTaskTableService", () => {
  let service: SystemTaskTableService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [NgxsModule.forRoot([SystemTaskTableState])],
      providers: [SystemTaskTableService]
    });
    service = TestBed.inject(SystemTaskTableService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
