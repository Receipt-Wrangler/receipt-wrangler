import { Component, OnInit } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { Group, GroupRole } from "../../open-api";
import { GroupUtil } from "../../utils";

@Component({
    selector: "app-group-details",
    templateUrl: "./group-details.component.html",
    styleUrl: "./group-details.component.scss",
    standalone: false
})
export class GroupDetailsComponent implements OnInit {
  public canEdit = false;

  public group!: Group;

  constructor(
    private groupUtil: GroupUtil,
    private activatedRoute: ActivatedRoute
  ) {
  }


  public ngOnInit(): void {
    this.group = this.activatedRoute.snapshot.data["group"];
    this.canEdit = this.groupUtil.hasGroupAccess(this.group.id, GroupRole.Owner, false);
  }
}
