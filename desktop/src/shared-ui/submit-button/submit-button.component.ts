import {
  Component,
  ViewEncapsulation,
  computed,
  input
} from '@angular/core';
import { Store } from '@ngxs/store';
import { FormMode } from 'src/enums/form-mode.enum';
import { LayoutState } from 'src/store/layout.state';
import { FormButtonComponent } from '../form-button/form-button.component';
import { ButtonComponent } from '../../button/button/button.component';

@Component({
    selector: 'app-submit-button',
    templateUrl: './submit-button.component.html',
    styleUrls: ['./submit-button.component.scss'],
    encapsulation: ViewEncapsulation.None,
    imports: [ButtonComponent]
})
export class SubmitButtonComponent extends FormButtonComponent {
  public readonly onlyIcon = input<boolean>(true);

  public readonly disableOnLoading = input<boolean>(false);

  public override readonly type = input<'button' | 'submit'>('submit');

  public formMode = FormMode;

  private showProgressBar = this.store.selectSignal(LayoutState.showProgressBar);

  public effectiveDisabled = computed(() =>
    this.disableOnLoading() && this.showProgressBar() ? true : this.disabled()
  );

  constructor(private store: Store) {
    super();
  }
}
