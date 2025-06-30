import { useParams, Link, useNavigate } from 'react-router-dom';
import { cars } from '@/data/cars';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft } from 'lucide-react';

const CarDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const car = cars.find((c) => c.id === id);

  if (!car) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Car not found</h1>
          <Link to="/">
            <Button>Return to Home</Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => navigate(-1)}
                className="flex items-center space-x-2"
              >
                <ArrowLeft className="w-4 h-4" />
                <span>Back</span>
              </Button>
              <div className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-blue-700 rounded-lg"></div>
                <h1 className="text-2xl font-bold text-gray-900">
                  DriveSelect
                </h1>
              </div>
            </div>
          </div>
        </div>
      </header>

      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Car Image */}
          <div className="space-y-4">
            <div className="relative overflow-hidden rounded-2xl shadow-2xl">
              <img
                src={car.image}
                alt={car.name}
                className="w-full h-96 lg:h-[500px] object-cover"
              />
              <div className="absolute top-6 left-6">
                <Badge className="bg-blue-600 text-white px-4 py-2 text-lg">
                  {car.type}
                </Badge>
              </div>
            </div>
          </div>

          {/* Car Details */}
          <div className="space-y-6">
            <div>
              <h1 className="text-4xl font-bold text-gray-900 mb-2">
                {car.brand} {car.name}
              </h1>
              <p className="text-xl text-gray-600 mb-4">
                {car.year} • {car.type}
              </p>
              <div className="flex items-center space-x-2 mb-6">
                <span className="text-3xl font-bold text-blue-600">
                  ${car.pricePerDay}
                </span>
                <span className="text-lg text-gray-500">per day</span>
              </div>
            </div>

            <Card className="border-0 shadow-lg">
              <CardHeader>
                <CardTitle className="text-xl">
                  Vehicle Specifications
                </CardTitle>
              </CardHeader>
              <CardContent className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <p className="text-sm text-gray-500">Transmission</p>
                  <p className="font-semibold">{car.transmission}</p>
                </div>
                <div className="space-y-2">
                  <p className="text-sm text-gray-500">Fuel Type</p>
                  <p className="font-semibold">{car.fuelType}</p>
                </div>
                <div className="space-y-2">
                  <p className="text-sm text-gray-500">Seats</p>
                  <p className="font-semibold">{car.seats} passengers</p>
                </div>
                <div className="space-y-2">
                  <p className="text-sm text-gray-500">Doors</p>
                  <p className="font-semibold">{car.doors} doors</p>
                </div>
              </CardContent>
            </Card>

            <Card className="border-0 shadow-lg">
              <CardHeader>
                <CardTitle className="text-xl">Features & Amenities</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  {car.features.map((feature, index) => (
                    <div key={index} className="flex items-center space-x-2">
                      <div className="w-2 h-2 bg-blue-600 rounded-full"></div>
                      <span className="text-gray-700">{feature}</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card className="border-0 shadow-lg">
              <CardHeader>
                <CardTitle className="text-xl">Description</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-700 leading-relaxed">
                  {car.description}
                </p>
              </CardContent>
            </Card>

            <div className="space-y-4">
              <a
                href={`http://cardealerform.local/quote/${car.id}`}
                className="block"
              >
                <Button
                  size="lg"
                  className="w-full bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-semibold py-4 text-lg rounded-xl transition-all duration-300 transform hover:scale-105"
                >
                  Get Quote & Book Now
                </Button>
              </a>
              <Link to="/">
                <Button
                  variant="outline"
                  size="lg"
                  className="w-full py-4 text-lg rounded-xl"
                >
                  Browse Other Cars
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CarDetail;
