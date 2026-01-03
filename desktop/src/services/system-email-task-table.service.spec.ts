import { TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { SystemEmailTaskTableState } from "../store/system-email-task-table.state";

import { SystemEmailTaskTableService } from "./system-email-task-table.service";

describe("SystemEmailTaskTableService", () => {
  let service: SystemEmailTaskTableService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        NgxsModule.forRoot([SystemEmailTaskTableState]),
      ],
      providers: [
        SystemEmailTaskTableService,
      ]
    });
    service = TestBed.inject(SystemEmailTaskTableService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
