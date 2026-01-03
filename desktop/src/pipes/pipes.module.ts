import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { CustomCurrencyPipe } from "./custom-currency.pipe";
import { CustomFieldTypePipe } from "./custom-field-type.pipe";
import { DurationPipe } from "./duration.pipe";
import { FormArrayLastPipe } from "./form-array-last.pipe";
import { FormGetPipe } from "./form-get.pipe";
import { GroupRolePipe } from "./group-role.pipe";
import { GroupPipe } from "./group.pipe";
import { ImagePipe } from "./image.pipe";
import { InputReadonlyPipe } from "./input-readonly.pipe";
import { MapGetPipe } from "./map-get.pipe";
import { MapKeyPipe } from "./map-key.pipe";
import { NamePipe } from "./name.pipe";
import { StatusPipe } from "./status.pipe";
import { UserPipe } from "./user.pipe";

@NgModule({
  declarations: [
    CustomCurrencyPipe,
    CustomFieldTypePipe,
    DurationPipe,
    FormArrayLastPipe,
    FormGetPipe,
    GroupPipe,
    GroupRolePipe,
    ImagePipe,
    InputReadonlyPipe,
    MapGetPipe,
    MapKeyPipe,
    NamePipe,
    StatusPipe,
    UserPipe,
  ],
  imports: [CommonModule],
  exports: [
    CustomCurrencyPipe,
    CustomFieldTypePipe,
    DurationPipe,
    FormArrayLastPipe,
    FormGetPipe,
    GroupPipe,
    GroupRolePipe,
    ImagePipe,
    InputReadonlyPipe,
    MapGetPipe,
    MapKeyPipe,
    NamePipe,
    StatusPipe,
    UserPipe,
  ],
})
export class PipesModule {}
