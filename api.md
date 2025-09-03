# Endpoints API

## Tipos

- `GET /api/types` - Listar todos los tipos de generadores
- `GET /api/types/{id}` - Obtener un tipo específico
- `POST /api/types` - Crear un nuevo tipo
- `PUT /api/types/{id}` - Actualizar un tipo existente
- `DELETE /api/types/{id}` - Eliminar un tipo
- `GET /api/types/{id}/generators` - Listar todos los generadores de un tipo específico

## Generadores

- `GET /api/generators` - Listar todos los generadores
- `GET /api/generators/{id}` - Obtener un generador específico
- `POST /api/generators` - Crear un nuevo generador
- `PUT /api/generators/{id}` - Actualizar un generador existente
- `DELETE /api/generators/{id}` - Eliminar un generador
- `GET /api/generators/{id}/productions` - Listar todos los registros de producción de un generador específico

### Producción por rango de fechas para un generador específico

- `GET /api/generators/{id}/production/{start}/{end}` - Listar registros de producción de un generador específico en un rango de fechas
- `GET /api/generators/{id}/production/{start}/` - Listar registros de producción de un generador específico desde una fecha en adelante

## Producción

- `GET /api/productions` - Listar todos los registros de producción
- `GET /api/productions/{id}` - Obtener un registro de producción específico
- `POST /api/productions` - Crear un nuevo registro de producción
- `PUT /api/productions/{id}` - Actualizar un registro de producción existente
- `DELETE /api/productions/{id}` - Eliminar un registro de producción
