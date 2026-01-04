import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { PromptTableState } from "../../store/prompt-table.state";
import { TableModule } from "../../table/table.module";

import { PromptTableComponent } from "./prompt-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("PromptTableComponent", () => {
  let component: PromptTableComponent;
  let fixture: ComponentFixture<PromptTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [PromptTableComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [TableModule, NgxsModule.forRoot([PromptTableState]), NoopAnimationsModule],
    providers: [{
            provide: ActivatedRoute,
            useValue: {
                snapshot: {
                    data: {
                        prompts: []
                    }
                }
            }
        }, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
})
      .compileComponents();

    fixture = TestBed.createComponent(PromptTableComponent);
    component = fixture.componentInstance;

  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
