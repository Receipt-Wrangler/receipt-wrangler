import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GroupAutocompleteComponent } from './group-autocomplete.component';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { NgxsModule } from '@ngxs/store';

describe('GroupAutocompleteComponent', () => {
  let component: GroupAutocompleteComponent;
  let fixture: ComponentFixture<GroupAutocompleteComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [GroupAutocompleteComponent],
      imports: [NgxsModule.forRoot([])],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
    });
    fixture = TestBed.createComponent(GroupAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
