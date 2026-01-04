import { provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { NgxsModule } from '@ngxs/store';
import { ApiModule } from '../open-api/api.module';
import { AppInitService } from './app-init.service';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

describe('AppInitService', () => {
  let service: AppInitService;

  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [NgxsModule.forRoot([]), ApiModule],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
    service = TestBed.inject(AppInitService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
