import { Component, Input, OnInit, input } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { BaseFormComponent } from "../../form";
import { FormConfig } from "../../interfaces";
import { Group, Prompt } from "../../open-api";
import { AutocomleteComponent } from "../../autocomplete/autocomlete/autocomlete.component";
import { FormGetPipe } from "../../pipes/form-get.pipe";
import { InputReadonlyPipe } from "../../pipes/input-readonly.pipe";

@Component({
    selector: "app-group-settings-ai",
    templateUrl: "./group-settings-ai.component.html",
    styleUrl: "./group-settings-ai.component.scss",
    imports: [AutocomleteComponent, FormGetPipe, InputReadonlyPipe]
})
export class GroupSettingsAiComponent extends BaseFormComponent implements OnInit {
  @Input() public override form: FormGroup = new FormGroup({});

  public readonly canEdit = input<boolean>(false);

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

    const canEdit = this.canEdit();
    if (!canEdit && this.group.groupSettings?.prompt) {
      prompts.push(this.group.groupSettings.prompt);
    }

    if (!canEdit && this.group.groupSettings?.fallbackPrompt) {
      prompts.push(this.group.groupSettings.fallbackPrompt);
    }

    if (canEdit) {
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
