# LazyMentor - Product Requirements Document (PRD)

## 1. Resumen Ejecutivo
**LazyMentor** es un agente mentor 100% en Markdown que vive dentro de cualquier agente de codificación IA (Cursor, Continue.dev / opencode, Claude Code, etc.).  
Su único propósito es ayudar al usuario que está migrando a **LazyVim** y aún no domina la navegación, los atajos y el flujo.

El proyecto consta de **dos partes**:
1. `lazymentor.md` → el prompt completo del agente (solo Markdown).
2. `lazymint-installer` → instalador TUI escrito **obligatoriamente en Go** que coloca el prompt en el agente elegido de forma cómoda y cross-platform.

## 2. Objetivos
- Acelerar el aprendizaje de LazyVim sin frustración.
- Priorizar **siempre** la configuración personal del usuario por encima de la documentación oficial.
- Respuestas **precisas, concisas y visuales** (tablas, listas, sin rodeos).
- **Cero riesgo**: el agente nunca genera código ni modifica archivos.

## 3. Requisitos Funcionales

### 3.1 Agente LazyMentor (`lazymentor.md`)
El archivo debe contener un **system prompt** completo con las siguientes reglas obligatorias:

**System Prompt Estructura Recomendada**

**Rol**
Eres LazyMentor, un mentor experto y paciente de LazyVim especializado SOLO en atajos y flujo de LazyVim . Hablas el idioma del usuario.

Reglas inquebrantables (siempre):
1. NUNCA generes ni una línea de código Lua, Vimscript, shell, etc.
2. NUNCA sugieras :w, :q, :e, ni ningún comando Ex que modifique archivos.
3. Respuestas ≤ 120 palabras + tabla si aplica.
4. Usa SIEMPRE tablas Markdown para listar keymaps.
5. Idioma: responde en el mismo idioma que la pregunta del usuario.

Prioridad de conocimiento (de mayor a menor):
1. Cualquier configuración que el usuario pegue en esta conversación (init.lua, keymaps.lua, etc.)
2. El popup de which-key que el usuario describa o pegue
   3. Documentación oficial de LazyVim[LazyVim](https://www.lazyvim.org/keymaps) como último recurso
Formato de respuesta preferido:
- Explicación corta (2–4 líneas)
- Tabla de keymaps relevantes cuando corresponda
- Pregunta de seguimiento si falta contexto: "¿Me puedes pegar tu mapping de <líder>e o el grupo que ves en which-key?"

**Regla de Oro #1 - Configuración del usuario**
- Siempre que el usuario mencione una tecla o comando, **pregunta primero** si quiere que uses su configuración personal.
- Si el usuario pega su `init.lua`, `lazy.lua` o cualquier parte de `~/.config/nvim`, **guárdala en memoria** (en el contexto de la conversación) y úsala como fuente de verdad #1. Nunca uses la documentación oficial si contradice su config.
- Si no tienes la config, pide educadamente: “¿Puedes pegarme la parte relevante de tu init.lua o el mapping que quieres usar?”

**Regla de Oro #2 - Formato de respuesta**
- Máximo 8 líneas de texto + tabla cuando corresponda.
- Usa **tablas Markdown** para atajos de teclado.
- Ejemplo requerido en el prompt:
  - Usuario: “cómo creo un archivo nuevo”
    Respuesta: “Espacio + e + a (según tu config). Si tu mapping es distinto, dime cuál usas.”
  - Usuario: “explica los buffers”
    Respuesta: explicación corta + tabla con los atajos esenciales.

**Regla de Oro #3 - Restricciones ABSOLUTAS**
- **NUNCA** generar código (ni un solo carácter).
- **NUNCA** sugerir editar, crear o modificar ningún archivo.
- **NUNCA** dar comandos `:w`, `:e`, `:qa`, etc. Solo explicar qué tecla usar.
- Si el usuario pide código o modificación → responder: “Lo siento, como LazyMentor no puedo generar código ni modificar archivos. ¿Quieres que te explique cómo hacerlo manualmente con atajos?”

**Ejemplos obligatorios que debe manejar**
- Crear archivo nuevo
- Abrir/cerrar buffers
- Navegación ( telescope, neo-tree, harpoon, etc.)
- Modo visual, selección, copiar/pegar
- Ventanas (splits)
- Cualquier atajo personalizado del usuario

### 3.2 Instalador TUI (`lazymint-installer` - Go)
**Tecnología obligatoria**: Go + TUI (recomendado: Charm Bracelet - bubbletea + lipgloss + progress).

**Funcionalidades del TUI (pantallas)**:
0. **Pre-flight checklist** (opcional, skipeable):
   - Confirmar que el usuario tiene Neovim abierto al lado
   - Confirmar que usa LazyVim o cualquier config de Neovim
   - Mensaje: "Esta herramienta funciona mejor si practicas los atajos en tu editor"
1. **Pantalla de bienvenida** + detección automática de OS (Linux / macOS / Windows).
2. **Selección del agente destino** (lista interactiva con flechas):
   - OpenCode
   - Claude Code
   - Otro (ruta manual)
3. **Detección automática de rutas**:
   - **OpenCode** → `~/.config/opencode/opencode.json`
     - El archivo `lazymentor.md` se embebe como entrada en `rules` o se copia como archivo referenciado.
   - **Claude Code** → `~/.claude/CLAUDE.md`
     - El contenido de `lazymentor.md` se agrega/crea como archivo global de instrucciones.
   - **Otro** → el usuario ingresa la ruta manualmente.
4. **Instalación en 1 clic**:
   - Copia `lazymentor.md` a la carpeta correcta del agente.
   - (Opcional) Agrega una entrada en el config del agente si es JSON.
   - Muestra barra de progreso + mensaje de éxito.
5. **Opciones extras**:
   - Instalar solo el MD en `~/lazymentor.md` (modo manual)
   - Actualizar versión existente
   - Desinstalar

**Flujo del instalador (origen de `lazymentor.md`)**:
El archivo `lazymentor.md` puede obtenerse de dos fuentes, en orden de prioridad:
1. **Archivo local junto al binario**: si `lazymentor.md` existe en el mismo directorio que el ejecutable, se usa ese.
2. **Embedded asset**: si no existe archivo local, se extrae del binario compilado (embebido en tiempo de compilación via Go's `embed`).
3. **Fallback de red**: si la extracción falla, se ofrece descargar la versión oficial desde el repositorio de GitHub.

```
┌─────────────────────────────────────────────────────────────┐
│  Usuario ejecuta lazymint-installer                         │
│  │                                                         │
│  ├─ ¿lazymentor.md existe junto al bin?                    │
│  │   ├─ SÍ → usar ese archivo (desarrollo local)           │
│  │   └─ NO →                                                   │
│  │       ├─ ¿Embedded extraction exitosa?                  │
│  │       │   ├─ SÍ → usar archivo embebido (release)      │
│  │       │   └─ NO →                                         │
│  │       │       └─ Ofrecer descarga desde GitHub         │
└─────────────────────────────────────────────────────────────┘
```

**Requisitos no funcionales del instalador**:
- 100% cross-platform (un solo binario).
- Sin dependencias externas (todo Go modules).
- Tamaño < 5 MB.
- Soporte para instalación con `go install` o release GitHub.
- Mensajes en English

**Manejo de errores y edge cases**:
El instalador debe manejar gracefully las siguientes situaciones:

| Escenario | Comportamiento esperado |
|-----------|------------------------|
| **Ruta del agente no existe** | Crear carpetas automáticamente o pedir confirmación al usuario |
| **Permisos denegados (EACCES)** | Mostrar error claro + sugerir `sudo` (Linux/mac) o ejecutar como administrador (Windows) |
| **Agente no detectado** | Ofrecer ingreso manual de ruta con validación |
| **Config JSON corrupto o mal formado** | Hacer backup automático antes de modificar (`.bak`), restaurar si la operación falla |
| **Ya está instalado** | Preguntar si quiere sobrescribir, actualizar, o cancelar |
| **Descarga de red falla** | Retry automático (3 intentos) + fallback al asset embebido con mensaje informativo |
| **Espacio en disco insuficiente** | Detectar anticipadamente y mostrar mensaje claro |
| **Instalación parcialmente completada** | Limpiar archivos creados antes de reportar error |

**Reglas de recuperación**:
- **Nunca dejar estado inconsistente**: si cualquier operación falla, rollback completo.
- **Backup automático**: antes de modificar cualquier archivo de configuración, crear `.lazymentor.backup.{timestamp}`.
- **Logging silencioso**: guardar un log de operaciones en `~/.lazymentor/install.log` para debugging.

## 4. Lo que NO puede hacer (hard constraints)
- El agente **nunca** genera código.
- El agente **nunca** toca archivos.
- El instalador **nunca** ejecuta comandos que modifiquen el nvim del usuario.
- No se permite usar Rust (instalador obligatorio en Go).

## 5. Arquitectura

lazymentor/
├── PRD.md
├── cmd/
│   └── installer/         ← main.go + CLI/TUI
├── internal/
│   ├── tui/               ← bubbletea screens (futuro)
│   ├── agents/            ← handlers por agente (OpenCode, Claude Code)
│   ├── config/            ← detección OS y rutas
│   └── embed/             ← lazymentor.md embebido (go:embed)
└── go.mod

## 6. Roadmap de desarrollo (prioridad)
1. Terminar y aprobar este PRD.
2. Codificar el instalador en Go (TUI).
3. Crear `lazymentor.md` (prompt definitivo).
4. Testing en los 3 agentes principales.
5. Release del binario.

## 7. Versionado y Releases

**Esquema de versionado**: [Semantic Versioning 2.0.0](https://semver.org/)

```
MAJOR.MINOR.PATCH
│     │     │
│     │     └─ Bug fixes, correcciones
│     └─ Nuevas funcionalidades (backwards compatible)
└─ Cambios rompe-compatibilidad
```

**Tags de Git**: `v1.0.0`, `v1.1.0`, `v2.0.0`, etc.

**Channels de release**:
| Channel | Tag | Frecuencia | Descripción |
|---------|-----|------------|-------------|
| **Stable** | `v*.*.*` | Cada 4-6 semanas | Versión probada y estable |
| **Beta** | `v*.*.*-beta.*` | Semanal | Nuevas features para testing |
| **Dev** | `v*.*.*-dev` | Continuous | Solo para desarrollo local |

**Artefactos por release**:
- `lazymint-installer-{os}-{arch}` — binario para cada SO/arquitectura
- `lazymentor.md` — archivo del prompt (embebido en binario + repo)

**Changelog**: Mantener `CHANGELOG.md` con formato [Keep a Changelog](https://keepachangelog.com/).

**Actualización del instalador**: Manual (v1.x)
- El instalador es un setup one-time; no se usa diariamente.
- Para actualizar: el usuario descarga la nueva versión desde GitHub releases manualmente.
- El comando `lazymint-installer update` NO existe en v1.x.
- Auto-update se considerará para v2.0 si hay demanda.

---
