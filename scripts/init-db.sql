-- Script de inicialización de la base de datos
-- Este script se ejecuta automáticamente cuando el contenedor se crea por primera vez

-- Crear extensiones útiles (opcional)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Puedes agregar aquí tablas iniciales, datos semilla, o configuraciones
-- Por ejemplo:
-- CREATE TABLE IF NOT EXISTS example (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- Mensaje de confirmación
DO $$
BEGIN
    RAISE NOTICE 'Base de datos DVRA inicializada correctamente';
END $$;
