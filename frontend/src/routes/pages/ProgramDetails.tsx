import { useParams } from 'react-router-dom';

export default function ProgramDetails() {
  const { id } = useParams<{ id: string }>();

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Program Details</h1>
      <p className="text-gray-600">Program ID: {id}</p>
    </div>
  );
}
