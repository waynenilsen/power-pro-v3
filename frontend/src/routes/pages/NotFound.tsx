import { Link } from 'react-router-dom';

export default function NotFound() {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">404 - Not Found</h1>
      <p className="mt-2 text-gray-600">The page you're looking for doesn't exist.</p>
      <Link to="/" className="mt-4 inline-block text-blue-600 hover:underline">
        Go back home
      </Link>
    </div>
  );
}
