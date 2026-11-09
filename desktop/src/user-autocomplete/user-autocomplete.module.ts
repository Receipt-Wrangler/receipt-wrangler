import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserAutocompleteComponent } from './user-autocomplete/user-autocomplete.component';
import { ReactiveFormsModule } from '@angular/forms';
import { AutocompleteModule } from 'src/autocomplete/autocomplete.module';
import { PipesModule } from 'src/pipes/pipes.module';

@NgModule({
  imports: [CommonModule, ReactiveFormsModule, AutocompleteModule, PipesModule, UserAutocompleteComponent],
  exports: [UserAutocompleteComponent],
})
export class UserAutocompleteModule {}
