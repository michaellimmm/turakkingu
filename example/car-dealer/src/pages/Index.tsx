import { Link } from 'react-router-dom';
import { cars } from '@/data/cars';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

const Index = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-blue-700 rounded-lg"></div>
              <h1 className="text-2xl font-bold text-gray-900">DriveSelect</h1>
            </div>
            <nav className="hidden md:flex items-center space-x-6">
              <a
                href="#cars"
                className="text-gray-600 hover:text-blue-600 transition-colors"
              >
                Cars
              </a>
              <a
                href="#about"
                className="text-gray-600 hover:text-blue-600 transition-colors"
              >
                About
              </a>
              <a
                href="#contact"
                className="text-gray-600 hover:text-blue-600 transition-colors"
              >
                Contact
              </a>
            </nav>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="relative py-20 overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-600 to-blue-800"></div>
        <div className="absolute inset-0 bg-black opacity-20"></div>
        <div className="relative container mx-auto px-4 text-center text-white">
          <h2 className="text-5xl md:text-6xl font-bold mb-6">
            Find Your Perfect
            <span className="block text-transparent bg-clip-text bg-gradient-to-r from-yellow-400 to-orange-500">
              Rental Car
            </span>
          </h2>
          <p className="text-xl md:text-2xl mb-8 text-blue-100 max-w-3xl mx-auto">
            Choose from our premium fleet of vehicles. From luxury sedans to
            rugged SUVs, we have the perfect car for every journey.
          </p>
          <Button
            size="lg"
            className="bg-white text-blue-600 hover:bg-blue-50 text-lg px-8 py-3 rounded-full font-semibold transition-all duration-300 transform hover:scale-105"
            onClick={() =>
              document
                .getElementById('cars')
                ?.scrollIntoView({ behavior: 'smooth' })
            }
          >
            Explore Our Fleet
          </Button>
        </div>
        <div className="absolute bottom-0 left-0 right-0 h-20 bg-gradient-to-t from-slate-50 to-transparent"></div>
      </section>

      {/* Cars Section */}
      <section id="cars" className="py-16">
        <div className="container mx-auto px-4">
          <div className="text-center mb-12">
            <h3 className="text-4xl font-bold text-gray-900 mb-4">
              Our Premium Fleet
            </h3>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Discover our carefully curated selection of vehicles, each
              maintained to the highest standards
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {cars.map((car) => (
              <Card
                key={car.id}
                className="group overflow-hidden hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2 bg-white border-0 shadow-lg"
              >
                <div className="relative overflow-hidden">
                  <img
                    src={car.image}
                    alt={car.name}
                    className="w-full h-56 object-cover group-hover:scale-110 transition-transform duration-500"
                  />
                  <div className="absolute top-4 left-4">
                    <Badge className="bg-blue-600 text-white px-3 py-1">
                      {car.type}
                    </Badge>
                  </div>
                  <div className="absolute top-4 right-4">
                    <Badge
                      variant="secondary"
                      className="bg-white text-gray-900 px-3 py-1 font-semibold"
                    >
                      ${car.pricePerDay}/day
                    </Badge>
                  </div>
                </div>

                <CardHeader className="pb-2">
                  <CardTitle className="text-xl font-bold text-gray-900">
                    {car.brand} {car.name}
                  </CardTitle>
                  <div className="flex items-center space-x-4 text-sm text-gray-600">
                    <span>{car.year}</span>
                    <span>•</span>
                    <span>{car.transmission}</span>
                    <span>•</span>
                    <span>{car.seats} seats</span>
                  </div>
                </CardHeader>

                <CardContent className="py-2">
                  <p className="text-gray-600 text-sm line-clamp-2">
                    {car.description}
                  </p>
                </CardContent>

                <CardFooter className="pt-2">
                  <Link to={`/car/${car.id}`} className="w-full">
                    <Button className="w-full bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-semibold py-2 rounded-lg transition-all duration-300">
                      View Details
                    </Button>
                  </Link>
                </CardFooter>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-12">
        <div className="container mx-auto px-4 text-center">
          <div className="flex items-center justify-center space-x-2 mb-6">
            <div className="w-8 h-8 bg-gradient-to-r from-blue-400 to-blue-500 rounded-lg"></div>
            <h3 className="text-2xl font-bold">DriveSelect</h3>
          </div>
          <p className="text-gray-400 mb-4">
            Your trusted partner for premium car rentals
          </p>
          <p className="text-gray-500 text-sm">
            © 2024 DriveSelect. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Index;
