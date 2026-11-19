import { Component, OnInit } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { Group, Permission } from "../../open-api";
import { AuthState } from "../../store";
import { GroupFormComponent } from "../group-form/group-form.component";

@Component({
    selector: "app-group-details",
    templateUrl: "./group-details.component.html",
    styleUrl: "./group-details.component.scss",
    imports: [GroupFormComponent]
})
export class GroupDetailsComponent implements OnInit {
  public canEdit = false;

  public group!: Group;

  constructor(
    private store: Store,
    private activatedRoute: ActivatedRoute
  ) {
  }


  public ngOnInit(): void {
    this.group = this.activatedRoute.snapshot.data["group"];
    this.canEdit = this.store.selectSnapshot(
      AuthState.hasGroupPermission(this.group.id, Permission.GroupUpdate)
    );
  }
}
