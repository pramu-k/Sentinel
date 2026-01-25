import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ServerDetails } from './server-details';

describe('ServerDetails', () => {
  let component: ServerDetails;
  let fixture: ComponentFixture<ServerDetails>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ServerDetails]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ServerDetails);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
