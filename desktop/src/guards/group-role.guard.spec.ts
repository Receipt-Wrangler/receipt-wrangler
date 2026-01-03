import { TestBed } from '@angular/core/testing';

import { GroupRoleGuard } from './group-role.guard';
import { RouterTestingModule } from '@angular/router/testing';
import { NgxsModule } from '@ngxs/store';

describe('GroupRoleGuard', () => {
  let guard: GroupRoleGuard;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [RouterTestingModule, NgxsModule.forRoot([])],
    });
    guard = TestBed.inject(GroupRoleGuard);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });
});
