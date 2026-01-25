import { TestBed } from '@angular/core/testing';

import { SentinelService } from './sentinel.service';

describe('SentinelService', () => {
  let service: SentinelService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SentinelService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
