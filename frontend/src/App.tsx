import { Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import SubjectList from './components/SubjectList'
import TopicList from './components/TopicList'
import FlashcardDeck from './components/FlashcardDeck'
import CreateFlashcard from './components/CreateFlashcard'
import AISubjectChat from './components/AISubjectChat'
import AIMarkingPractice from './components/AIMarkingPractice'
import SubjectProgressPage from './components/SubjectProgress'
import TextToFlashcards from './components/TextToFlashcards'
import Login from './components/Login'
import Signup from './components/Signup'
import CompleteSignup from './components/CompleteSignup'
import TutorList from './components/TutorList'
import TutorDetailPage from './components/TutorDetail'
import StudentEssays from './components/StudentEssays'
import TutorEssays from './components/TutorEssays'
import EssayDetailPage from './components/EssayDetail'
import { useAuth } from './auth/AuthContext'

function RequireAuth({ children }: { children: JSX.Element }) {
  const { user, loading } = useAuth()
  if (loading) {
    return <p style={{ color: 'var(--muted)' }}>Loading…</p>
  }
  if (!user) {
    return <Navigate to="/login" replace />
  }
  return children
}

export default function App() {
  return (
    <Routes>
      {/* Public auth routes */}
      <Route
        path="/login"
        element={
          <Layout>
            <Login />
          </Layout>
        }
      />
      <Route
        path="/signup"
        element={
          <Layout>
            <Signup />
          </Layout>
        }
      />
      <Route
        path="/signup/complete"
        element={
          <Layout>
            <CompleteSignup />
          </Layout>
        }
      />

      {/* Protected app routes */}
      <Route
        path="/"
        element={
          <RequireAuth>
            <Layout>
              <SubjectList />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/subjects/:subjectId/topics"
        element={
          <RequireAuth>
            <Layout>
              <TopicList />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/subjects/:subjectId/practice"
        element={
          <RequireAuth>
            <Layout>
              <AIMarkingPractice />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/subjects/:subjectId/progress"
        element={
          <RequireAuth>
            <Layout>
              <SubjectProgressPage />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/topics/:topicId/cards"
        element={
          <RequireAuth>
            <Layout>
              <FlashcardDeck />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/topics/:topicId/cards/new"
        element={
          <RequireAuth>
            <Layout>
              <CreateFlashcard />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/topics/:topicId/ai/text"
        element={
          <RequireAuth>
            <Layout>
              <TextToFlashcards />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/ai/subjects"
        element={
          <RequireAuth>
            <Layout>
              <AISubjectChat />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/tutors"
        element={
          <RequireAuth>
            <Layout>
              <TutorList />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/tutors/:id"
        element={
          <RequireAuth>
            <Layout>
              <TutorDetailPage />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/student/essays"
        element={
          <RequireAuth>
            <Layout>
              <StudentEssays />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/tutor/essays"
        element={
          <RequireAuth>
            <Layout>
              <TutorEssays />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/essays/:id"
        element={
          <RequireAuth>
            <Layout>
              <EssayDetailPage />
            </Layout>
          </RequireAuth>
        }
      />
    </Routes>
  )
}
