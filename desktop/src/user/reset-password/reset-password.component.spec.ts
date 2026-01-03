import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NgxsModule } from "@ngxs/store";
import { PipesModule } from "src/pipes/pipes.module";
import { ApiModule } from "../../open-api";
import { ResetPasswordComponent } from "./reset-password.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ResetPasswordComponent", () => {
  let component: ResetPasswordComponent;
  let fixture: ComponentFixture<ResetPasswordComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [ResetPasswordComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        PipesModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([]),
        PipesModule,
        ReactiveFormsModule],
    providers: [{ provide: MatDialogRef, useValue: {} }, provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
}).compileComponents();

    fixture = TestBed.createComponent(ResetPasswordComponent);
    component = fixture.componentInstance;
    component.user = { id: "" } as any;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });
});
