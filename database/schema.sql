-- ============================================================
-- BLISS HOSPITAL - SMART APPOINTMENT & QUEUE MANAGEMENT SYSTEM
-- Database Schema (PostgreSQL via Supabase)
-- Author: Cherise Akinyi Otieno | J17/0967/2022
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- DEPARTMENTS
-- ============================================================
CREATE TABLE departments (
  department_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name            VARCHAR(255) NOT NULL,
  description     TEXT,
  created_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE users (
  user_id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  full_name       VARCHAR(255) NOT NULL,
  email           VARCHAR(255) UNIQUE NOT NULL,
  phone           VARCHAR(20),
  password_hash   TEXT NOT NULL,
  role            VARCHAR(20) DEFAULT 'patient'
    CHECK (role IN ('patient', 'receptionist', 'doctor', 'admin')),
  is_active       BOOLEAN DEFAULT TRUE,
  created_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- BRANCHES (✅ FIX ADDED)
-- ============================================================
CREATE TABLE branches (
  branch_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  branch_name VARCHAR(255) NOT NULL,
  created_at  TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- DOCTORS
-- ============================================================
CREATE TABLE doctors (
  doctor_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id         UUID REFERENCES users(user_id) ON DELETE SET NULL,
  department_id   UUID REFERENCES departments(department_id) ON DELETE SET NULL,
  full_name       VARCHAR(255) NOT NULL,
  specialty       VARCHAR(255),
  bio             TEXT,
  available       BOOLEAN DEFAULT TRUE,
  created_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- AVAILABILITY SLOTS
-- ============================================================
CREATE TABLE availability_slots (
  slot_id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  doctor_id       UUID REFERENCES doctors(doctor_id) ON DELETE CASCADE,
  slot_date       DATE NOT NULL,
  start_time      TIME NOT NULL,
  end_time        TIME NOT NULL,
  is_booked       BOOLEAN DEFAULT FALSE,
  created_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- APPOINTMENTS
-- ============================================================
CREATE TABLE appointments (
  appointment_id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id         UUID REFERENCES users(user_id) ON DELETE CASCADE,
  doctor_id       UUID REFERENCES doctors(doctor_id) ON DELETE SET NULL,
  department_id   UUID REFERENCES departments(department_id) ON DELETE SET NULL,
  branch_id       UUID REFERENCES branches(branch_id) ON DELETE SET NULL,
  slot_id         UUID REFERENCES availability_slots(slot_id) ON DELETE SET NULL,
  appointment_date DATE NOT NULL,
  appointment_time TIME NOT NULL,
  reason          TEXT,
  status          VARCHAR(20) DEFAULT 'pending'
    CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed', 'no_show')),
  notes           TEXT,
  created_at      TIMESTAMP DEFAULT NOW(),
  updated_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- QUEUE
-- ============================================================
CREATE TABLE queue (
  queue_id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  appointment_id  UUID REFERENCES appointments(appointment_id) ON DELETE CASCADE,
  user_id         UUID REFERENCES users(user_id) ON DELETE CASCADE,
  doctor_id       UUID REFERENCES doctors(doctor_id) ON DELETE SET NULL,
  queue_number    INT NOT NULL,
  position        INT NOT NULL,
  status          VARCHAR(20) DEFAULT 'waiting'
    CHECK (status IN ('waiting', 'called', 'serving', 'done', 'skipped')),
  checked_in_at   TIMESTAMP,
  called_at       TIMESTAMP,
  served_at       TIMESTAMP,
  created_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- NOTIFICATIONS
-- ============================================================
CREATE TABLE notifications (
  notification_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id         UUID REFERENCES users(user_id) ON DELETE CASCADE,
  appointment_id  UUID REFERENCES appointments(appointment_id) ON DELETE SET NULL,
  type            VARCHAR(30)
    CHECK (type IN ('confirmation', 'reminder', 'delay', 'cancellation', 'queue_call')),
  channel         VARCHAR(10)
    CHECK (channel IN ('sms', 'email')),
  message         TEXT NOT NULL,
  sent_at         TIMESTAMP,
  status          VARCHAR(20) DEFAULT 'pending'
    CHECK (status IN ('pending', 'sent', 'failed')),
  created_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- AUDIT LOG
-- ============================================================
CREATE TABLE audit_log (
  log_id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id         UUID REFERENCES users(user_id) ON DELETE SET NULL,
  action          VARCHAR(100) NOT NULL,
  table_affected  VARCHAR(100),
  record_id       UUID,
  details         TEXT,
  ip_address      VARCHAR(50),
  created_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- INDEXES
-- ============================================================
CREATE INDEX idx_appointments_user   ON appointments(user_id);
CREATE INDEX idx_appointments_doctor ON appointments(doctor_id);
CREATE INDEX idx_appointments_date   ON appointments(appointment_date);
CREATE INDEX idx_appointments_status ON appointments(status);
CREATE INDEX idx_queue_doctor        ON queue(doctor_id);
CREATE INDEX idx_queue_status        ON queue(status);
CREATE INDEX idx_slots_doctor_date   ON availability_slots(doctor_id, slot_date);
CREATE INDEX idx_notifications_user  ON notifications(user_id);

-- ============================================================
-- TRIGGERS
-- ============================================================
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER appointments_updated_at
BEFORE UPDATE ON appointments
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================================
-- SEED DATA (DEPARTMENTS)
-- ============================================================
INSERT INTO departments (name, description) VALUES
('General Medicine',  'Primary care and general health consultations'),
('Cardiology',        'Heart and cardiovascular system specialists'),
('Dermatology',       'Skin, hair and nail conditions'),
('Pediatrics',        'Healthcare for infants, children and adolescents'),
('Gynecology',        'Female reproductive health'),
('Orthopedics',       'Bone, joint and muscle conditions'),
('Ophthalmology',     'Eye care and vision health'),
('Dental',            'Oral health and dental procedures'),
('ENT',               'Ear, nose and throat specialists'),
('Mental Health',     'Psychiatric and psychological services');

-- ============================================================
-- SEED DATA (BRANCHES)
-- ============================================================
INSERT INTO branches (branch_name) VALUES
('Moi Avenue - Nairobi'),
('Eastleigh (BBS Mall) - Nairobi'),
('Langata - Nairobi'),
('Githurai - Nairobi'),
('Pipeline - Nairobi'),
('Laiboni Centre (Lenana Road) - Nairobi'),
('Thika (KRA Building)'),
('Murang’a'),
('Embu'),
('Meru'),
('Mombasa (Kizingo)'),
('Mtwapa - Mombasa'),
('Nakuru');