import { Routes } from '@angular/router';
import { UserListComponent } from './user-list/user-list.component';

export const userRoutes: Routes = [
  {
    path: '',
    component: UserListComponent,
  },
];
