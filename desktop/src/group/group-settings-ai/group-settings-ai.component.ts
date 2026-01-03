import { Component, Input, OnInit } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { BaseFormComponent } from "../../form";
import { FormConfig } from "../../interfaces";
import { Group, Prompt } from "../../open-api";

@Component({
    selector: "app-group-settings-ai",
    templateUrl: "./group-settings-ai.component.html",
    styleUrl: "./group-settings-ai.component.scss",
    standalone: false
})
export class GroupSettingsAiComponent extends BaseFormComponent implements OnInit {
  @Input() public override form: FormGroup = new FormGroup({});

  @Input() public canEdit: boolean = false;

  @Input() public override formConfig!: FormConfig;

  public group!: Group;

  public prompts: Prompt[] = [];

  constructor(
    private activatedRoute: ActivatedRoute,
  ) {
    super();
  }

  public ngOnInit(): void {
    this.group = this.activatedRoute.snapshot.data["group"];
    this.setPrompts();
    this.initForm();
  }

  private setPrompts(): void {
    let prompts: Prompt[] = [];

    if (!this.canEdit && this.group.groupSettings?.prompt) {
      prompts.push(this.group.groupSettings.prompt);
    }

    if (!this.canEdit && this.group.groupSettings?.fallbackPrompt) {
      prompts.push(this.group.groupSettings.fallbackPrompt);
    }

    if (this.canEdit) {
      prompts = this.activatedRoute.snapshot.data["prompts"];
    }

    this.prompts = prompts;
  }

  private initForm(): void {
    this.setInitialValues();
  }

  private setInitialValues(): void {
    this.emitFormCommand(
      {
        path: "promptId",
        command: "patchValue",
        payload: this.group.groupSettings?.promptId ?? null
      }
    );
    this.emitFormCommand(
      {
        path: "fallbackPromptId",
        command: "patchValue",
        payload: this.group.groupSettings?.fallbackPromptId ?? null
      }
    );
  }

  public promptDisplayWith(id: number): string {
    return this.prompts.find((prompt) => prompt.id === id)?.name ?? "";
  }
}
