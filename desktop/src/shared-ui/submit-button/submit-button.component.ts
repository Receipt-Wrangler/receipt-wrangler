import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  ViewEncapsulation,
} from '@angular/core';
import { Store } from '@ngxs/store';
import { Subscription, startWith, tap } from 'rxjs';
import { FormMode } from 'src/enums/form-mode.enum';
import { LayoutState } from 'src/store/layout.state';
import { FormButtonComponent } from '../form-button/form-button.component';

@Component({
    selector: 'app-submit-button',
    templateUrl: './submit-button.component.html',
    styleUrls: ['./submit-button.component.scss'],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class SubmitButtonComponent
  extends FormButtonComponent
  implements OnChanges
{
  @Input() public onlyIcon: boolean = true;

  @Input() public disableOnLoading: boolean = false;

  @Input() public override type: 'button' | 'submit' = 'submit';

  public formMode = FormMode;

  public disabledSubscription!: Subscription;

  public originalDisabledState: boolean = this.disabled;

  constructor(private store: Store) {
    super();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    this.originalDisabledState = this.disabled;

    if (changes['disableOnLoading']?.currentValue || changes['disabled']) {
      this.disabledSubscription = this.store
        .select(LayoutState.showProgressBar)
        .pipe(
          startWith(this.store.selectSnapshot(LayoutState.showProgressBar)),
          tap((isShown) => {
            if (isShown) {
              this.disabled = true;
            } else {
              this.disabled = this.originalDisabledState;
            }
          })
        )
        .subscribe();
    } else {
      if (this.disabledSubscription) {
        this.disabledSubscription.unsubscribe();
      }
    }
  }
}
