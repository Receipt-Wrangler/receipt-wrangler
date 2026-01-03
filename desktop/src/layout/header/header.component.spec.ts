import { provideHttpClientTesting } from "@angular/common/http/testing";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NgxsModule, Store } from "@ngxs/store";
import { ToggleIsSidebarOpen } from "src/store/layout.state.actions";
import { ApiModule } from "../../open-api";
import { HeaderComponent } from "./header.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("HeaderComponent", () => {
  let component: HeaderComponent;
  let fixture: ComponentFixture<HeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [HeaderComponent],
    imports: [ApiModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([])],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
}).compileComponents();

    fixture = TestBed.createComponent(HeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should toggle sidebar", () => {
    const store = jest.spyOn(TestBed.inject(Store), "dispatch");
    component.toggleSidebar();

    expect(store).toHaveBeenCalledWith(new ToggleIsSidebarOpen());
  });
});
