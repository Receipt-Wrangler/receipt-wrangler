import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { NgxsModule } from "@ngxs/store";
import { ImportType } from "../../open-api";
import { PipesModule } from "../../pipes";

import { ImportFormComponent } from "./import-form.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ImportFormComponent", () => {
  let component: ImportFormComponent;
  let fixture: ComponentFixture<ImportFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [ImportFormComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ReactiveFormsModule,
        MatDialogModule,
        PipesModule,
        NgxsModule.forRoot([])],
    providers: [
        { provide: MatDialogRef, useValue: {} },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
    ]
})
      .compileComponents();

    fixture = TestBed.createComponent(ImportFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      importType: ImportType.ImportConfig
    });
  });
});
