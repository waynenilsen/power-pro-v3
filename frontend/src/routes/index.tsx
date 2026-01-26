import { createBrowserRouter } from 'react-router-dom';
import { Layout } from '../components/layout';
import Home from './pages/Home';
import Login from './pages/Login';
import Programs from './pages/Programs';
import ProgramDetails from './pages/ProgramDetails';
import Workout from './pages/Workout';
import History from './pages/History';
import Profile from './pages/Profile';
import NotFound from './pages/NotFound';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <Home />,
      },
      {
        path: 'programs',
        element: <Programs />,
      },
      {
        path: 'programs/:id',
        element: <ProgramDetails />,
      },
      {
        path: 'workout',
        element: <Workout />,
      },
      {
        path: 'history',
        element: <History />,
      },
      {
        path: 'profile',
        element: <Profile />,
      },
    ],
  },
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '*',
    element: <NotFound />,
  },
]);
