# Smart Task Planner - AI-Powered Task Management System

## Overview
Smart Task Planner is an intelligent task management platform that leverages AI to help users break down goals into actionable tasks, track progress, analyze risks, and receive personalized feedback. Built with Go and MongoDB, it features natural language processing for intuitive task management and AI-driven insights.

**API Testing**: All endpoints have been tested using Postman with comprehensive screenshots provided below.

## Key Features

- **AI-Powered Task Generation**: Automatically breaks down goals into detailed, actionable tasks
- **Smart Scheduling**: Intelligent deadline assignment avoiding conflicts and risky dates
- **Risk Analysis**: Identifies tasks at risk of missing deadlines
- **Progress Tracking**: Real-time completion percentage with AI-generated motivational feedback
- **Natural Language Commands**: Interact using plain English queries
- **Task Refinement**: Break down complex tasks into subtasks using AI
- **Auto-Rescheduling**: Bulk deadline adjustments when falling behind
- **Google OAuth Integration**: Secure authentication with JWT tokens
- **Hierarchical Task Structure**: Unlimited nesting of subtasks

## Technology Stack

- **Backend**: Go 1.21+ with Gin framework
- **Database**: MongoDB
- **AI Integration**: OpenAI GPT-4 & Google Gemini
- **Authentication**: JWT + Google OAuth 2.0
- **Architecture**: MCP (Model Context Protocol) for AI tool orchestration

---

## API Documentation

### Authentication Endpoints

#### Register User
```
POST /api/auth/register
```
Create a new user account with email and password.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepass123"
}
```

**Response:**
```json
{
  "token": "jwt-token-here",
  "name": "John Doe",
  "email": "john@example.com"
}
```

#### Login User
```
POST /api/auth/login
```
Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepass123"
}
```

**Response:**
```json
{
  "token": "jwt-token-here",
  "name": "John Doe",
  "email": "john@example.com"
}
```

#### Google OAuth Login
```
GET /api/auth/google/login
```
Redirects to Google OAuth consent screen.

#### Google OAuth Callback
```
GET /api/auth/google/callback
```
Handles OAuth callback and returns JWT token.

---

### Plan Management Endpoints

> **Note**: All plan endpoints require `Authorization: Bearer <token>` header

#### Generate Draft Plan (AI)
```
POST /api/plan/draft
```
Generate an AI-powered task breakdown without saving to database.

**Request Body:**
```json
{
  "goal": "Learn machine learning in 3 months"
}
```

**Response:**
```json
{
  "goal": "Learn machine learning in 3 months",
  "tasks": [
    {
      "title": "Set up Python development environment",
      "description": "Install Python 3.x, Jupyter Notebook, and essential ML libraries",
      "status": "Pending",
      "deadline": "2025-10-20T00:00:00Z"
    },
    {
      "title": "Learn Python fundamentals",
      "description": "Master Python basics: variables, loops, functions, OOP",
      "status": "Pending",
      "deadline": "2025-10-25T00:00:00Z"
    }
    // ... 8 more tasks
  ]
}
```

#### Confirm and Save Plan
```
POST /api/plan/confirm
```
Save the generated plan to database.

**Request Body:**
```json
{
  "goal": "Learn machine learning in 3 months"
}
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "goal": "Learn machine learning in 3 months",
  "tasks": [/* task array */]
}
```

#### Get All Plans
```
GET /api/plan/
```
Retrieve all plans for the authenticated user.

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "user_id": "user-id-here",
    "goal": "Learn machine learning in 3 months",
    "tasks": [/* task array */]
  }
]
```

#### Refine Task into Subtasks
```
POST /api/plan/refine-task
```
Break down a complex task into detailed subtasks using AI.

**Request Body:**
```json
{
  "task": {
    "title": "Build a neural network",
    "description": "Create and train a basic neural network"
  }
}
```

**Response:**
```json
{
  "subtasks": [
    {
      "title": "Design network architecture",
      "description": "Define input, hidden, and output layers",
      "status": "Pending",
      "deadline": "2025-10-22T00:00:00Z"
    },
    {
      "title": "Implement forward propagation",
      "description": "Code the forward pass through the network",
      "status": "Pending",
      "deadline": "2025-10-24T00:00:00Z"
    }
    // ... more subtasks
  ]
}
```

#### Update Task Status
```
POST /api/plan/update-task-status
```
Mark a task as completed or update its status.

**Request Body:**
```json
{
  "task_id": "507f1f77bcf86cd799439011",
  "status": "Completed"
}
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Set up Python development environment",
  "status": "Completed",
  "deadline": "2025-10-20T00:00:00Z"
}
```

#### Get Task Details
```
GET /api/plan/task-details?task_id=507f1f77bcf86cd799439011
```
Retrieve detailed information about a specific task.

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Learn Python fundamentals",
  "description": "Master Python basics",
  "status": "In Progress",
  "deadline": "2025-10-25T00:00:00Z",
  "sub_tasks": []
}
```

#### Add Subtasks to Task
```
POST /api/plan/add-subtasks
```
Add AI-generated or manual subtasks to an existing task.

**Request Body:**
```json
{
  "task_id": "507f1f77bcf86cd799439011",
  "subtasks": [
    {
      "title": "Learn variables and data types",
      "description": "Understand int, float, string, list, dict",
      "deadline": "2025-10-22T00:00:00Z"
    }
  ]
}
```

---

### Natural Language Command Endpoints

> **Note**: Requires `Authorization: Bearer <token>` header

#### Process Natural Language Command
```
POST /api/command/
```
Send natural language queries to interact with your tasks.

**Request Body:**
```json
{
  "message": "Show me my progress on machine learning"
}
```

**Supported Commands:**

1. **Progress & Feedback**
```json
{
  "message": "What's my progress on learning Python?"
}
```

**Response:**
```json
{
  "tool": "get_user_progress",
  "result": {
    "user_id": "user-id",
    "goal": "Learn machine learning in 3 months",
    "completion_percentage": 35,
    "total_tasks": 10,
    "completed_tasks": 3
  },
  "feedback": {
    "tone": "You're making progress! 📈",
    "message": "Great work! You're 35% done with 'Learn machine learning'. 3 tasks completed, 7 to go!",
    "suggestion": "You're doing well, but let's pick up the pace a bit. Try to complete more tasks this week to reach 50%!",
    "progress_summary": {
      "goal": "Learn machine learning in 3 months",
      "completion_percentage": 35,
      "completed_tasks": 3,
      "remaining_tasks": 7,
      "total_tasks": 10
    }
  }
}
```

2. **Risk Analysis**
```json
{
  "message": "What tasks are at risk?"
}
```

**Response:**
```json
{
  "tool": "analyze_risks",
  "result": {
    "user_id": "user-id",
    "risks": [
      {
        "goal": "Learn machine learning",
        "task_name": "Complete linear regression project",
        "deadline": "2025-10-18",
        "days_left": 2
      },
      {
        "goal": "Learn machine learning",
        "task_name": "Study decision trees",
        "deadline": "2025-10-19",
        "days_left": 3
      }
    ],
    "count": 2,
    "threshold_days": 3
  }
}
```

3. **Reschedule Plan**
```json
{
  "message": "I'm behind schedule by 5 days on my machine learning goal"
}
```

**Response:**
```json
{
  "tool": "reschedule_plan",
  "result": {
    "message": "All tasks for goal 'Learn machine learning' shifted by 5 days",
    "goal_id": "507f1f77bcf86cd799439011",
    "tasks": [/* updated task array with new deadlines */]
  }
}
```

4. **Alternative Plans**
```json
{
  "message": "Can I finish faster?"
}
```

**Response:**
```json
{
  "tool": "generate_alternative_plans",
  "result": {
    "goal_id": "507f1f77bcf86cd799439011",
    "options": [
      {
        "type": "speed",
        "description": "Focus on completing tasks faster, may reduce quality."
      },
      {
        "type": "balance",
        "description": "Balanced approach between speed and quality."
      },
      {
        "type": "quality",
        "description": "Focus on doing tasks with highest quality, may take longer."
      }
    ]
  }
}
```

5. **General Queries** (AI-Powered)
```json
{
  "message": "What should I focus on this week?"
}
```

**Response:**
```json
{
  "tool": "handle_general_query",
  "result": {
    "response": "Based on your current progress, you should prioritize completing 'Learn Python fundamentals' which is due in 2 days. You're making good progress at 35% completion. Try to wrap up your in-progress tasks before starting new ones!",
    "context_used": true
  }
}
```

---

### Health Check Endpoints

#### Server Health
```
GET /health
```
Check if the server and database are operational.

**Response:**
```json
{
  "status": "healthy",
  "message": "Server is running",
  "database": "connected"
}
```

#### Root Endpoint
```
GET /
```
Get API information.

**Response:**
```json
{
  "message": "Smart Task Planner API",
  "version": "1.0.0"
}
```

---

## Data Models

### User
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "John Doe",
  "email": "john@example.com",
  "password": "hashed-password"
}
```

### Plan
```json
{
  "id": "507f1f77bcf86cd799439011",
  "user_id": "user-id-here",
  "goal": "Learn machine learning in 3 months",
  "tasks": [
    {
      "id": "507f191e810c19729de860ea",
      "title": "Set up Python environment",
      "description": "Install Python and required libraries",
      "status": "Completed",
      "deadline": "2025-10-20T00:00:00Z",
      "sub_tasks": []
    }
  ]
}
```

### Task (Recursive)
```json
{
  "id": "507f191e810c19729de860ea",
  "title": "Build neural network",
  "description": "Create and train a basic neural network",
  "status": "Pending",
  "deadline": "2025-11-15T00:00:00Z",
  "sub_tasks": [
    {
      "id": "507f191e810c19729de860eb",
      "title": "Design architecture",
      "description": "Define layers and activation functions",
      "status": "Pending",
      "deadline": "2025-11-10T00:00:00Z",
      "sub_tasks": []
    }
  ]
}
```

---

## Postman Testing Screenshots

### Authentication
![User Registration](<assets/register.png>)
*User registration with email and password*

![User Login](<assets/login.png>)
*Login endpoint returning JWT token*

![Google OAuth Flow](<screenshots/auth-google-oauth.png>)
*Google OAuth authentication*

### Plan Management
![Generate Draft Plan](<assets/plandraft1.png>)
*AI-generated task breakdown for a goal*

![Confirm Plan](<assets/planconfirm.png>)
*Saving the generated plan to database*

![Get All Plans](<assets/getplans.png>)
*Retrieving all user plans*

![Refine Task](<assets/refinetask.png>)
*Breaking down a task into subtasks*

![Add SubTasks](<assets/addsubtasks.png>)
*adding subtasks to the task*

![Update Task Status](<assets/updatetaskstatus.png>)
*Marking a task as completed*

### Natural Language Commands
![Progress Query](<assets/command1.png>)
![Progress Query](<assets/command2.png>)
*Asking about goal progress with AI feedback*

![Risk Analysis](<assets/command3.png>)
*Checking tasks at risk of missing deadlines*

![Reschedule Plan](<assets/command5.png>)
*Bulk rescheduling tasks when behind*

![General Query](<assets/command4.png>)
*AI-powered response to custom questions*

### Health Checks
![Health Endpoint](<assets/health.png>)
*Server and database health status*

---

## Installation & Setup

### Prerequisites
- Go 1.21 or higher
- MongoDB 6.0+
- OpenAI API Key
- Google Cloud OAuth credentials (optional)

### Backend Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd smart-task-planner
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Set up environment variables**
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```env
# Server Configuration
PORT=8080

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=smart_task_planner

# Google OAuth (Optional)
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# JWT Secret
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# AI API Keys
OPENAI_API_KEY=sk-your-openai-api-key
GEMINI_API_KEY=your-gemini-api-key-optional
```

4. **Run the application**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Docker Setup (Optional)

```bash
# Build and run with Docker Compose
docker-compose up --build -d

# View logs
docker-compose logs -f

# Stop containers
docker-compose down
```

---

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `PORT` | Server port | Yes | `8080` |
| `MONGODB_URI` | MongoDB connection string | Yes | - |
| `MONGODB_DATABASE` | Database name | Yes | - |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | No | - |
| `GOOGLE_CLIENT_SECRET` | Google OAuth secret | No | - |
| `GOOGLE_REDIRECT_URL` | OAuth callback URL | No | - |
| `JWT_SECRET` | Secret for JWT signing | Yes | - |
| `OPENAI_API_KEY` | OpenAI API key | Yes | - |
| `GEMINI_API_KEY` | Google Gemini API key | No | - |

---

## MCP (Model Context Protocol) - The AI Brain

### What is MCP?

**Model Context Protocol (MCP)** is the intelligent orchestration layer that transforms Smart Task Planner from a traditional CRUD application into an **AI-native assistant**. It acts as a unified interface between users, AI models, and your task management system.

Think of MCP as your project's **"smart dispatcher"** - it interprets what users want, decides which specialized tool should handle it, executes that tool with proper context, and returns intelligent responses.

### Why MCP Exists

#### **Problem Without MCP**
Traditional task management apps require users to:
- Remember specific API endpoints (`/api/plan/reschedule`, `/api/plan/analyze-risks`)
- Know exact request formats
- Make multiple API calls for complex operations
- Understand technical terminology

**Example**: To check progress and get feedback, users would need:
```bash
# Call 1: Get progress
POST /api/plan/progress

# Call 2: Calculate stats manually

# Call 3: Get feedback
POST /api/plan/feedback
```

#### **Solution With MCP**
Users simply ask questions in plain English:
```bash
POST /api/command/
{
  "message": "How am I doing on my machine learning goal?"
}
```

MCP automatically:
1. Interprets the intent (progress query)
2. Fetches relevant data
3. Calculates statistics
4. Generates AI-powered feedback
5. Returns a complete, contextual response

### How MCP Benefits Users

| Without MCP | With MCP |
|-------------|----------|
| "I need to call `/api/plan/reschedule`" | "I'm behind by 3 days" |
| "Which endpoint shows risks?" | "What tasks are at risk?" |
| "How do I format a task refinement request?" | "Break down this task for me" |
| Multiple API calls for one question | Single natural language query |
| Technical knowledge required | Conversational interface |

### MCP Architecture in Our Project

```
┌─────────────────────────────────────────────────────────────┐
│                        USER INPUT                            │
│  "What's my progress?" / "I'm behind" / "Show risks"        │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                   COMMAND HANDLER                            │
│              POST /api/command/                              │
│         (Receives natural language input)                    │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                  MCP EXECUTOR (Brain)                        │
│                  internal/mcp/executor.go                    │
│                                                              │
│  RunTool(tool_name, params, repo) → Routes to tools         │
└────────────────────────────┬────────────────────────────────┘
                             │
                ┌────────────┴────────────┐
                │                         │
                ▼                         ▼
┌──────────────────────────┐  ┌──────────────────────────┐
│   INTERPRET MESSAGE      │  │    SPECIFIC MCP TOOLS    │
│  tools_phase3.go         │  │                          │
│                          │  │  • create_task_plan      │
│  Keyword Detection:      │  │  • analyze_risks         │
│  • "progress" →          │  │  • get_user_progress     │
│    get_user_progress     │  │  • reschedule_plan       │
│  • "behind" →            │  │  • refine_task           │
│    reschedule_plan       │  │  • provide_feedback      │
│  • "risk" →              │  │  • handle_general_query  │
│    analyze_risks         │  │                          │
│  • No match →            │  └──────────┬───────────────┘
│    handle_general_query  │             │
└──────────┬───────────────┘             │
           │                             │
           └──────────┬──────────────────┘
                      │
        ┌─────────────┴─────────────┐
        │                           │
        ▼                           ▼
┌─────────────────┐        ┌─────────────────┐
│  AI SERVICES    │        │   DATABASE      │
│                 │        │                 │
│  • OpenAI API   │        │  • MongoDB      │
│  • Gemini API   │        │  • User Plans   │
│  • GPT-4        │        │  • Task Data    │
└────────┬────────┘        └────────┬────────┘
         │                          │
         └──────────┬───────────────┘
                    │
                    ▼
         ┌─────────────────────┐
         │  STRUCTURED RESULT  │
         │  + AI INSIGHTS      │
         │  + RECOMMENDATIONS  │
         └─────────────────────┘
```

### MCP Tool Categories

#### 🤖 **AI-Powered Tools** (Leverage OpenAI/Gemini)

**1. `create_task_plan`** (`create_task_plan.go`)
```go
// Purpose: Generate comprehensive task breakdown from high-level goals
// AI Model: OpenAI GPT-4
// Intelligence: 
//   - Analyzes goal complexity
//   - Creates 10+ actionable tasks
//   - Assigns smart deadlines avoiding conflicts
//   - Provides detailed descriptions

Goal: "Learn machine learning in 3 months"
↓
AI generates:
[
  {"title": "Set up Python environment", "deadline": "2025-10-20"},
  {"title": "Learn NumPy and Pandas", "deadline": "2025-10-27"},
  {"title": "Study linear regression", "deadline": "2025-11-03"},
  // ... 7 more tasks with even distribution
]
```

**2. `refine_task`** (`refine_task.go`)
```go
// Purpose: Break complex tasks into manageable subtasks
// AI Model: OpenAI GPT-4
// Intelligence:
//   - Understands task complexity
//   - Creates logical subtask hierarchy
//   - Maintains deadline coherence

Task: "Build a neural network"
↓
AI breaks down into:
[
  {"title": "Design architecture", "deadline": "2025-11-10"},
  {"title": "Implement forward propagation", "deadline": "2025-11-12"},
  {"title": "Code backpropagation", "deadline": "2025-11-15"},
  {"title": "Train and validate model", "deadline": "2025-11-18"}
]
```

**3. `handle_general_query`** (`tools_phase3.go`)
```go
// Purpose: Answer custom questions with full context awareness
// AI Model: OpenAI GPT-4
// Intelligence:
//   - Fetches all user plans
//   - Calculates progress for each goal
//   - Identifies upcoming deadlines
//   - Generates contextual advice

User: "What should I prioritize this week?"
↓
MCP builds context:
"User's plans:
 1. Learn ML (35% complete, Python basics due in 2 days)
 2. Build portfolio (60% complete, no urgent tasks)"
↓
AI response:
"Focus on completing Python basics first since it's due in 2 days. 
You're making good progress on ML - maintain this momentum!"
```

**4. `FindGoalByAI`** (Helper in `repository`)
```go
// Purpose: Match user's vague references to actual goals
// AI Model: OpenAI GPT-4o-mini
// Intelligence:
//   - Semantic matching of user intent to goals
//   - Handles typos and variations

User: "shift my ML plan"
↓
AI matches to: "Learn machine learning in 3 months"
(even though user didn't type exact goal name)
```

#### 📊 **Analytics Tools** (Database + Logic)

**5. `analyze_risks`** (`tools_phase3.go`)
```go
// Purpose: Identify tasks at risk of missing deadlines
// Logic: Recursively scans all tasks & subtasks
// Output: Sorted list by urgency (least days remaining first)

Algorithm:
1. Fetch all user plans
2. For each task and subtask:
   - Calculate days_left = deadline - today
   - If days_left ≤ threshold (default 3):
     → Add to risk list
3. Sort by days_left (ascending)
4. Return with goal context

Result: [
  {goal: "Learn ML", task: "Complete project", days_left: 1},
  {goal: "Learn ML", task: "Study decision trees", days_left: 2}
]
```

**6. `get_user_progress`** (`tools_phase3.go`)
```go
// Purpose: Calculate completion statistics
// Logic: Recursive task counting with status filtering

Algorithm:
1. Find relevant goal (AI-powered if user mentions specific goal)
2. Recursively count:
   - Total tasks (including all nested subtasks)
   - Completed tasks (status == "Completed")
3. Calculate: progress = (completed / total) * 100
4. Return structured statistics

Result: {
  goal: "Learn machine learning",
  completion_percentage: 35,
  total_tasks: 10,
  completed_tasks: 3,
  remaining_tasks: 7
}
```

#### 🛠️ **Action Tools** (Data Modification)

**7. `reschedule_plan`** (`tools_phase3.go`)
```go
// Purpose: Bulk deadline adjustments when users fall behind
// Intelligence: Regex extraction + AI goal matching

Flow:
1. Extract delay days from message:
   "I'm behind by 5 days" → extracts "5"
2. Use AI to find which goal user means
3. Loop through all tasks:
   - Add delay days to each deadline
   - Preserve deadline relationships
4. Update MongoDB atomically

Result: All tasks shifted by specified days, maintaining relative spacing
```

**8. `update_task_status`** (`update_task_status.go`)
```go
// Purpose: Mark tasks as completed/in-progress
// Database: Direct MongoDB update using positional operator

Update pattern:
filter: {"tasks._id": task_id}
update: {"$set": {"tasks.$.status": "Completed"}}
        ↑ Positional operator updates matched array element
```

#### 💬 **Feedback Tools** (Rule-Based Intelligence)

**9. `provide_feedback`** (`tools_phase3.go`)
```go
// Purpose: Generate motivational messages based on progress
// Logic: Rule-based personalization

Algorithm:
if completion == 0%:
  tone = "Let's get started! 🚀"
  message = "The best time to begin is now!"
  
else if completion < 25%:
  tone = "Good start! 💪"
  suggestion = "Complete 2-3 tasks this week"
  
else if completion < 50%:
  tone = "You're making progress! 📈"
  suggestion = "Pick up the pace to reach 50%"
  
// ... progressive encouragement
  
else if completion == 100%:
  tone = "Goal achieved! 🎉"
  message = "Time to set a new goal!"

Response includes:
- Personalized tone
- Contextual message
- Actionable suggestion
- Progress summary
```

### How MCP Processes Requests

#### **Example 1: Simple Keyword Match**

```
User Input: "Show me tasks at risk"
↓
1. interpret_user_message() detects keyword "risk"
   ↓
2. Returns: {
     tool: "analyze_risks",
     params: {user_id: "123"}
   }
   ↓
3. MCP Executor calls analyze_risks()
   ↓
4. Queries MongoDB for all user plans
   ↓
5. Calculates days_left for each task
   ↓
6. Filters tasks with days_left ≤ 3
   ↓
7. Sorts by urgency
   ↓
8. Returns: {
     risks: [
       {goal: "Learn ML", task: "Project", days_left: 1},
       {goal: "Learn ML", task: "Study DT", days_left: 2}
     ],
     count: 2
   }
```

#### **Example 2: AI-Powered Query with Chaining**

```
User Input: "How am I doing on my machine learning goal?"
↓
1. interpret_user_message() detects "progress"
   ↓
2. Returns: {
     tool: "get_user_progress",
     params: {user_id: "123", message: "..."},
     needs_chaining: true  // ← Triggers automatic feedback
   }
   ↓
3. MCP Executor calls get_user_progress()
   ↓
4. Uses AI to match "machine learning goal" to exact goal
   ↓
5. Recursively counts tasks: 10 total, 3 completed
   ↓
6. Calculates: 30% complete
   ↓
7. MCP sees needs_chaining = true
   ↓
8. Automatically calls provide_feedback()
   ↓
9. Rule engine generates message:
   "You're making progress! 📈 
    You're 30% done with 'Learn machine learning'. 
    3 tasks completed, 7 to go!
    Keep building momentum!"
   ↓
10. Returns combined result with both stats and feedback
```

#### **Example 3: Fallback to General AI**

```
User Input: "What's the best way to learn faster?"
↓
1. interpret_user_message() checks keywords:
   - No "progress", "risk", "behind", "faster" match
   ↓
2. Falls through to default case
   ↓
3. Returns: {
     tool: "handle_general_query",
     params: {user_id: "123", message: "..."}
   }
   ↓
4. MCP Executor calls handle_general_query()
   ↓
5. Fetches all user plans from MongoDB
   ↓
6. Builds rich context:
   "User's plans:
    1. Learn ML (30% complete)
       - Upcoming: Python basics (due in 2 days)
    2. Build portfolio (60% complete)
       - No urgent tasks"
   ↓
7. Constructs AI prompt:
   "User asked: 'What's the best way to learn faster?'
    
    Context: [user's plan data]
    
    Provide helpful advice (2-3 sentences max)"
   ↓
8. Calls OpenAI API
   ↓
9. AI generates contextual response:
   "Based on your ML progress, focus on completing one task 
    fully before starting new ones. Your Python basics task 
    is due soon - prioritize that. Break complex tasks into 
    smaller chunks to maintain momentum!"
   ↓
10. Returns AI response with context_used: true flag
```

### Why MCP is Powerful

#### **1. Natural Language Interface**
Users don't need API documentation:
```
❌ Without MCP:
POST /api/plan/507f1f77bcf86cd799439011/tasks/507f191e810c19729de860ea/reschedule
Body: {"delay_days": 3, "cascade": true}

✅ With MCP:
POST /api/command/
Body: {"message": "I'm 3 days behind on my ML goal"}
```

#### **2. Context Awareness**
MCP has full access to user data:
```go
// MCP can automatically:
- Fetch all user plans
- Calculate progress across goals
- Identify deadline conflicts
- Understand task relationships
- Provide personalized advice
```

#### **3. Tool Chaining**
MCP can automatically trigger related tools:
```
User asks for progress
  ↓
MCP calls get_user_progress
  ↓
Automatically chains to provide_feedback
  ↓
Single response with complete information
```

#### **4. Intelligent Routing**
MCP decides which tool to use:
```go
// User says: "I'm behind"
// MCP routes to: reschedule_plan

// User says: "What should I do?"
// MCP routes to: handle_general_query (AI-powered)

// User says: "Am I on track?"
// MCP routes to: get_user_progress + provide_feedback
```

#### **5. Extensibility**
Adding new features is trivial:
```go
// Add new tool in executor.go:
case "send_reminder_email":
    return sendReminderEmail(params, repo)

// That's it! Now users can say:
// "Remind me about my deadline tomorrow"
```

### MCP vs Traditional REST API

| Aspect | Traditional REST | MCP Approach |
|--------|------------------|--------------|
| **Learning Curve** | High (must know endpoints) | Low (natural language) |
| **Request Format** | Structured JSON | Plain English |
| **Error Handling** | HTTP codes + messages | Conversational error explanations |
| **Multi-Step Operations** | Multiple API calls | Single command |
| **Discoverability** | Need documentation | Ask questions naturally |
| **Context Retention** | Stateless (must send context each time) | Context-aware |
| **AI Integration** | Separate endpoints | Unified through MCP |
| **User Experience** | Technical | Conversational |

### Real-World Impact

**Scenario**: User falls behind on a goal

**Without MCP (6 steps):**
1. GET /api/plan/ to find goal ID
2. GET /api/plan/{id} to see all tasks
3. Manually calculate new deadlines
4. For each task: PUT /api/plan/{id}/tasks/{task_id}
5. Repeat for all subtasks
6. Verify changes with another GET

**With MCP (1 step):**
```bash
POST /api/command/
{"message": "I'm 5 days behind on my ML goal"}

# MCP automatically:
# - Finds the right goal using AI
# - Updates all tasks and subtasks
# - Maintains deadline relationships
# - Returns confirmation
```

**Time saved**: ~10 minutes → 5 seconds

### Technical Excellence

#### **Smart Prompt Engineering**
```go
// In handle_general_query:
prompt := fmt.Sprintf(`You are a helpful task planning assistant.
A user asked: "%s"

Here's what you know about the user:
%s

Provide a helpful, conversational response that:
1. Directly answers their question
2. Is encouraging and supportive
3. Suggests actionable next steps
4. Keep it concise (2-3 sentences max)

Response:`, message, context)
```

This ensures AI responses are:
- Relevant (uses actual user data)
- Actionable (suggests next steps)
- Concise (respects user time)
- Encouraging (maintains motivation)

#### **Robust Error Handling**
```go
// All MCP tools have graceful fallbacks
aiResponse, err := CallOpenAIAPI(prompt)
if err != nil {
    return map[string]interface{}{
        "response": "I understand your question, but I'm having 
                     trouble generating a detailed response. 
                     Try asking about your progress, risks, 
                     or rescheduling tasks!",
    }, nil
}
```

#### **Efficient Database Queries**
```go
// analyze_risks uses indexed queries
plans, err := repo.GetAllByUser(userID) // Indexed on user_id

// Single query fetches all data
// In-memory processing for speed
```

### Future MCP Enhancements

The MCP architecture makes these additions trivial:

1. **Email Reminders**: `send_deadline_reminder`
2. **Calendar Sync**: `sync_to_google_calendar`
3. **Voice Commands**: `interpret_voice_message`
4. **Team Collaboration**: `share_task_with_team`
5. **Smart Suggestions**: `suggest_next_task`
6. **Time Tracking**: `log_time_spent`
7. **Goal Templates**: `apply_goal_template`

Each is just a new tool in the MCP ecosystem!

---

## AI Features Explained

### 1. Task Generation (via MCP)
The system uses OpenAI GPT-4 through the `create_task_plan` MCP tool to:
- Break down high-level goals into 10+ actionable tasks
- Assign realistic deadlines distributed across the goal timeline
- Avoid scheduling conflicts with existing risky tasks
- Provide detailed descriptions for each task

### 2. Smart Scheduling (via MCP)
The `create_task_plan` tool includes intelligent scheduling:
- Analyzes existing task deadlines
- Identifies "risky dates" (tasks due within 3 days)
- Automatically adjusts new task deadlines to avoid conflicts
- Ensures even distribution of workload

### 3. Progress Feedback (via MCP)
The `provide_feedback` tool generates personalized motivational messages based on:
- Completion percentage (0%, <25%, <50%, <75%, <90%, 100%)
- Number of remaining tasks
- Time since last update
- Overall goal complexity

### 4. Natural Language Processing (via MCP)
The `interpret_user_message` function routes queries to appropriate MCP tools:
- **Keywords detected**: "behind", "risk", "faster", "progress", "feedback"
- **Fallback**: AI-powered `handle_general_query` tool for custom questions
- **Context-aware**: Uses user's plan data to generate relevant responses

### 5. Risk Analysis (via MCP)
The `analyze_risks` tool:
- Scans all tasks and subtasks recursively
- Identifies tasks with deadlines ≤ 3 days away
- Sorts by urgency (least days remaining first)
- Provides actionable alerts

---

## Architecture Overview

```
cmd/server/
├── main.go                 # Application entry point

internal/
├── config/                 # Configuration management
├── database/              # MongoDB connection
├── middleware/            # JWT authentication
├── modules/
│   ├── auth/             # Authentication (JWT + OAuth)
│   ├── command/          # Natural language command handler
│   ├── plan/             # Plan CRUD operations
│   └── mcp/              # AI tool orchestration
│       ├── executor.go           # Tool router
│       ├── create_task_plan.go   # AI task generation
│       ├── tools_phase3.go       # Risk, progress, feedback
│       ├── refine_task.go        # Subtask generation
│       └── openai.go             # OpenAI API client
```

### Request Flow
```
User Request → Gin Router → JWT Middleware → Handler → Service → Repository → MongoDB
                                                  ↓
                                            MCP Executor → AI Tools → OpenAI API
```

---

## Performance & Best Practices

### Database Optimization
- **Indexes**: Created on `user_id`, `tasks._id`, and `email` fields
- **Query Optimization**: Uses BSON filters for efficient lookups
- **Connection Pooling**: Reuses MongoDB connections

### AI Optimization
- **Prompt Engineering**: Concise, structured prompts for consistent responses
- **Error Handling**: Graceful fallbacks when AI fails
- **JSON Parsing**: Robust handling of AI-generated JSON with markdown stripping

### Security
- **JWT Tokens**: Secure authentication with expiration
- **Password Hashing**: Bcrypt for password storage
- **Environment Variables**: Sensitive keys never hardcoded
- **OAuth 2.0**: Secure Google authentication flow

### Rate Limiting
- Consider implementing rate limiting for AI endpoints
- Suggested: 10 requests/minute per user for AI-heavy operations

---

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "details": "Additional context"
}
```

Common HTTP Status Codes:
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request body or parameters
- `401 Unauthorized`: Missing or invalid JWT token
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

---

## Testing with Postman

### Setting Up Postman

1. **Import Collection** (if provided)
   - File → Import → Select collection JSON

2. **Configure Environment Variables**
   - Create a new environment
   - Add variables:
     - `base_url`: `http://localhost:8080`
     - `token`: (will be set after login)

3. **Authentication Flow**
   - Register/Login → Copy JWT token
   - Set `Authorization` header: `Bearer {{token}}`
   - All subsequent requests will use this token

### Sample Postman Collection Structure

```
Smart Task Planner/
├── Authentication/
│   ├── Register
│   ├── Login
│   └── Google OAuth
├── Plans/
│   ├── Generate Draft
│   ├── Confirm Plan
│   ├── Get All Plans
│   ├── Refine Task
│   ├── Update Status
│   └── Get Task Details
├── Commands/
│   ├── Progress Query
│   ├── Risk Analysis
│   ├── Reschedule
│   └── General Query
└── Health/
    └── Health Check
```

---

## Troubleshooting

### Common Issues

**MongoDB Connection Failed**
```bash
# Check MongoDB is running
mongosh

# Verify connection string in .env
MONGODB_URI=mongodb://localhost:27017
```

**OpenAI API Errors**
```bash
# Verify API key is valid
echo $OPENAI_API_KEY

# Check API usage limits
# Visit: https://platform.openai.com/usage
```

**JWT Token Invalid**
```bash
# Ensure JWT_SECRET is set
# Token expires after 24 hours by default
# Re-login to get a fresh token
```

**Google OAuth Fails**
```bash
# Verify redirect URL matches Google Console
# Check GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET
# Ensure OAuth consent screen is configured
```

---

## Future Enhancements

- [ ] WebSocket support for real-time updates
- [ ] Email notifications for task deadlines
- [ ] Calendar integration (Google Calendar, Outlook)
- [ ] Mobile app (React Native)
- [ ] Team collaboration features
- [ ] Advanced analytics dashboard
- [ ] Voice command support
- [ ] Task templates and presets
- [ ] Recurring task scheduling
- [ ] File attachments for tasks

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Contact & Support

For questions, issues, or suggestions:
- **Email**: support@smarttaskplanner.com
- **GitHub Issues**: [Create an issue](https://github.com/yourusername/smart-task-planner/issues)
- **Documentation**: [Full API docs](https://docs.smarttaskplanner.com)

---

## Acknowledgments

- OpenAI for GPT-4 API
- Google for Gemini API
- MongoDB team for excellent documentation
- Gin framework contributors

---

**Built with ❤️ using Go, MongoDB, and AI**