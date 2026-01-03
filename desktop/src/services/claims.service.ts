import { Injectable } from '@angular/core';
import { Store } from '@ngxs/store';
import { Observable, switchMap, take } from 'rxjs';
import { UserService } from '../open-api/api/user.service';
import { SetAuthState } from '../store/auth.state.actions';

@Injectable({
  providedIn: 'root',
})
export class ClaimsService {
  constructor(private store: Store, private userService: UserService) {}

  public getAndSetClaimsForLoggedInUser(): Observable<void> {
    return this.userService.getUserClaims().pipe(
      take(1),
      switchMap((claims) => this.store.dispatch(new SetAuthState(claims)))
    );
  }
}
