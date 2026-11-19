import { Component, ElementRef, input, viewChild } from "@angular/core";
import { MatAutocompleteTrigger, MatAutocomplete, MatOption } from "@angular/material/autocomplete";
import { BaseInputComponent } from "../../base-input";
import { InputInterface } from "../../input";
import { MatFormField, MatLabel, MatInput, MatHint, MatError } from "@angular/material/input";
import { CdkTextareaAutosize } from "@angular/cdk/text-field";
import { ReactiveFormsModule } from "@angular/forms";
import { AsyncPipe } from "@angular/common";

@Component({
    selector: "app-textarea",
    templateUrl: "./textarea.component.html",
    styleUrls: ["./textarea.component.scss"],
    imports: [MatFormField, MatLabel, CdkTextareaAutosize, MatInput, ReactiveFormsModule, MatAutocompleteTrigger, MatHint, MatAutocomplete, MatOption, MatError, AsyncPipe]
})
export class TextareaComponent
  extends BaseInputComponent
  implements InputInterface {
  public readonly matAutocompleteTrigger = viewChild.required(MatAutocompleteTrigger);

  public readonly textarea = viewChild.required<ElementRef<HTMLTextAreaElement>>("nativeTextarea");

  public readonly options = input<string[]>([]);

  public readonly trigger = input<string>("");

  public filteredOptions: string[] = [];

  public lastKnownSelection: number = -1;

  public validEndCharacters = [" ", "\n", undefined];

  public onOptionSelected(): void {
    const value = this.inputFormControl.value;
    // Calculate insertion index for autocomplete selection
    let insertionIndex = value.slice(this.lastKnownSelection).search(/[\s\n]|$/);
    insertionIndex = insertionIndex === -1 ? value.length : this.lastKnownSelection + insertionIndex;
    this.matAutocompleteTrigger().closePanel();
    this.textarea().nativeElement.selectionEnd = insertionIndex + 1;
  }

  public onSelectionChange(): void {
    this.lastKnownSelection = this.textarea().nativeElement.selectionStart;
    // Extract word at cursor position for triggering autocomplete
    const currentWordDetails = this.getTriggerWordFromIndex(this.lastKnownSelection - 1);
    if (currentWordDetails.word !== null) {  // Check if a trigger word exists
      this.matAutocompleteTrigger().openPanel();
      this.filterOptions(currentWordDetails.word);
    } else {
      this.matAutocompleteTrigger().closePanel();
    }
  }

  private filterOptions(currentWord: string): void {
    if (currentWord === "") {  // Check if the extracted word is empty, indicating just the trigger character
      this.filteredOptions = this.options();  // Show all options if only the trigger character is typed
    } else {
      this.filteredOptions = this.options().filter(option =>
        option.toLowerCase().startsWith(currentWord.toLowerCase())
      );
    }
  }

  private getTriggerWordFromIndex(index: number): { word: string | null, triggerIndex: number } {
    const preText = this.inputFormControl.value.substring(0, index + 1);
    // Regex to capture word following the trigger character, allowing for empty follow-up
    const match = preText.match(new RegExp(`\\${this.trigger()}([^${this.validEndCharacters.join("")}]*)$`));
    if (match && match.index !== undefined) {
      return { word: match[1], triggerIndex: match.index };
    }
    return { word: null, triggerIndex: -1 };  // Return null if no trigger character is found
  }

  public getOptionValue(option: string): string {
    const insertionIndex = this.textarea().nativeElement.selectionEnd;
    const value = this.inputFormControl.value;
    const triggerWordDetails = this.getTriggerWordFromIndex(insertionIndex - 1);
    return value.slice(0, triggerWordDetails.triggerIndex) + this.trigger() + option + value.slice(insertionIndex) + " ";
  }
}
