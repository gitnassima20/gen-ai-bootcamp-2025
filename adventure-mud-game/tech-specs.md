# Adventure MUD Game

## Business Goal

**Goal**: A language immersion text-based adventure game using Go + Fiber

**Focus**: Simple MUD (Multi-User Dungeon) mechanics with predefined rooms, items, and NPCs

**Language Learning**: Each response includes one Japanese vocabulary word


## Technical Requirements

- Backend: Go + Fiber
- AI/NPCs: AWS Bedrock (nova-light) maybe
- Frontend: Minimal Go Templates + HTMX
- State Management: In-Memory (Go Maps)
- Add the 2D to enable interaction with screen

## Database Schema

for later

## API Endpoints

## AI Features

- [x] NPC Dialogues: Powered by AWS Bedrock (nova-light)
- [ ] AI Room Lore Expansion: When using Look command
- [ ] AI Puzzle Hints
- [ ] Background Audio Matching based On Room
- [ ] Character Builder
- [ ] Vocab Quest (example:
🌸 Quest: Find three items whose names include the kanji for "light" (光))
- [ ] Guardrails for responses(NPC stays in character, Prevent NSFW, religious, or political references, Limit Hints... )
