CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE bookings (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   title VARCHAR(255) NOT NULL,
   start_date TIMESTAMPTZ NOT NULL,
   end_date TIMESTAMPTZ NOT NULL,
   is_repeat BOOLEAN NOT NULL DEFAULT false,
   user_id UUID NOT NULL,
   office_id UUID NOT NULL,
   room_id UUID NOT NULL,
   CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id),
   CONSTRAINT fk_office_id FOREIGN KEY (office_id) REFERENCES offices(id),
   CONSTRAINT fk_room_id FOREIGN KEY (room_id) REFERENCES rooms(id),
   slug VARCHAR(255) NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP
)