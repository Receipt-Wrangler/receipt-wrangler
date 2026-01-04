import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";
import { NgxsModule } from "@ngxs/store";
import { of } from "rxjs";
import { PromptService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";

import { PromptFormComponent } from "./prompt-form.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("PromptFormComponent", () => {
  let component: PromptFormComponent;
  let fixture: ComponentFixture<PromptFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [PromptFormComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [SharedUiModule, PipesModule, ReactiveFormsModule, NgxsModule.forRoot([])],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: {
                snapshot: {
                    data: {
                        prompt: {},
                        formConfig: {}
                    }
                }
            }
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
    ]
})
      .compileComponents();

    fixture = TestBed.createComponent(PromptFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form without data", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      name: null,
      description: null,
      prompt: null,
    });
  });

  it("should init form with data", () => {
    const activatedRoute = TestBed.inject(ActivatedRoute);
    const prompt = {
      name: "ChatGPT",
      description: "ChatGPT prompt",
      prompt: "do the magic!"
    } as any;
    activatedRoute.snapshot.data["prompt"] = prompt;

    component.ngOnInit();

    expect(component.form.value).toEqual({
      name: prompt.name,
      description: prompt.description,
      prompt: prompt.prompt,
    });
  });

  it("prompt should be invalid due to bad variables", () => {
    component.ngOnInit();
    component?.form?.get("prompt")?.setValue("do the @magic!");

    expect(component.form.get("prompt")?.errors).toEqual({ invalidVariables: { value: true } });
  });

  it("prompt should be valid", () => {
    component.ngOnInit();
    component?.form?.get("prompt")?.setValue("please use @categories, and @tags");

    expect(component.form.get("prompt")?.errors).toEqual(null);
  });

  it("should create prompt", () => {
    const promptService = TestBed.inject(PromptService);
    const router = TestBed.inject(Router);
    const snackbarService = TestBed.inject(SnackbarService);

    const createPromptSpy = jest.spyOn(promptService, "createPrompt").mockReturnValue(of({} as any));
    const routerNavigateSpy = jest.spyOn(router, "navigate");
    const snackbarServiceSpy = jest.spyOn(snackbarService, "success");

    component.ngOnInit();
    component.originalPrompt = undefined;
    component.form.setValue({
      name: "ChatGPT",
      description: "ChatGPT prompt",
      prompt: "do the magic!"
    });

    component.submit();

    expect(createPromptSpy).toHaveBeenCalled();
    expect(routerNavigateSpy).toHaveBeenCalled();
    expect(snackbarServiceSpy).toHaveBeenCalled();
  });

  it("should update prompt", () => {
    const promptService = TestBed.inject(PromptService);
    const router = TestBed.inject(Router);
    const snackbarService = TestBed.inject(SnackbarService);

    const updatePromptSpy = jest.spyOn(promptService, "updatePromptById").mockReturnValue(of({} as any));
    const routerNavigateSpy = jest.spyOn(router, "navigate");
    const snackbarServiceSpy = jest.spyOn(snackbarService, "success");

    component.ngOnInit();
    component.originalPrompt = {
      id: 1,
      name: "ChatGPT",
      description: "ChatGPT prompt",
      prompt: "do the magic!"
    } as any;
    component.form.setValue({
      name: "ChatGPT",
      description: "ChatGPT prompt",
      prompt: "do the magic!"
    });

    component.submit();

    expect(updatePromptSpy).toHaveBeenCalled();
    expect(routerNavigateSpy).toHaveBeenCalled();
    expect(snackbarServiceSpy).toHaveBeenCalled();
  });
});
