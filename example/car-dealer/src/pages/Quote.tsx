import { useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { cars } from '@/data/cars';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { ArrowLeft } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

const Quote = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { toast } = useToast();
  const car = cars.find((c) => c.id === id);

  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
    pickupDate: '',
    returnDate: '',
    pickupLocation: '',
    driverLicense: '',
    additionalRequests: '',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

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

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.firstName.trim())
      newErrors.firstName = 'First name is required';
    if (!formData.lastName.trim()) newErrors.lastName = 'Last name is required';
    if (!formData.email.trim()) newErrors.email = 'Email is required';
    if (!/\S+@\S+\.\S+/.test(formData.email))
      newErrors.email = 'Email is invalid';
    if (!formData.phone.trim()) newErrors.phone = 'Phone number is required';
    if (!formData.pickupDate) newErrors.pickupDate = 'Pickup date is required';
    if (!formData.returnDate) newErrors.returnDate = 'Return date is required';
    if (!formData.pickupLocation.trim())
      newErrors.pickupLocation = 'Pickup location is required';
    if (!formData.driverLicense.trim())
      newErrors.driverLicense = 'Driver license number is required';

    // Date validation
    if (formData.pickupDate && formData.returnDate) {
      const pickup = new Date(formData.pickupDate);
      const returnDate = new Date(formData.returnDate);
      const today = new Date();
      today.setHours(0, 0, 0, 0);

      if (pickup < today) {
        newErrors.pickupDate = 'Pickup date cannot be in the past';
      }
      if (returnDate <= pickup) {
        newErrors.returnDate = 'Return date must be after pickup date';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const calculateDays = () => {
    if (formData.pickupDate && formData.returnDate) {
      const pickup = new Date(formData.pickupDate);
      const returnDate = new Date(formData.returnDate);
      const diffTime = returnDate.getTime() - pickup.getTime();
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
      return Math.max(1, diffDays);
    }
    return 1;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (validateForm()) {
      // Simulate form submission
      toast({
        title: 'Quote Request Submitted!',
        description: 'Your rental request has been processed successfully.',
      });

      // Redirect to thank you page after a short delay
      setTimeout(() => {
        navigate('/thank-you');
      }, 1500);
    } else {
      toast({
        title: 'Please fix the errors',
        description: 'Check the form for any missing or invalid information.',
        variant: 'destructive',
      });
    }
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: '' }));
    }
  };

  const days = calculateDays();
  const totalPrice = days * car.pricePerDay;

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
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-4">
              Request Your Quote
            </h1>
            <p className="text-xl text-gray-600">
              Fill out the form below to get started with your rental
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Form */}
            <div className="lg:col-span-2">
              <Card className="border-0 shadow-lg">
                <CardHeader>
                  <CardTitle className="text-2xl">Rental Information</CardTitle>
                </CardHeader>
                <CardContent>
                  <form onSubmit={handleSubmit} className="space-y-6">
                    {/* Personal Information */}
                    <div>
                      <h3 className="text-lg font-semibold mb-4 text-gray-900">
                        Personal Information
                      </h3>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <Label htmlFor="firstName">First Name *</Label>
                          <Input
                            id="firstName"
                            value={formData.firstName}
                            onChange={(e) =>
                              handleInputChange('firstName', e.target.value)
                            }
                            className={errors.firstName ? 'border-red-500' : ''}
                          />
                          {errors.firstName && (
                            <p className="text-red-500 text-sm mt-1">
                              {errors.firstName}
                            </p>
                          )}
                        </div>
                        <div>
                          <Label htmlFor="lastName">Last Name *</Label>
                          <Input
                            id="lastName"
                            value={formData.lastName}
                            onChange={(e) =>
                              handleInputChange('lastName', e.target.value)
                            }
                            className={errors.lastName ? 'border-red-500' : ''}
                          />
                          {errors.lastName && (
                            <p className="text-red-500 text-sm mt-1">
                              {errors.lastName}
                            </p>
                          )}
                        </div>
                        <div>
                          <Label htmlFor="email">Email *</Label>
                          <Input
                            id="email"
                            type="email"
                            value={formData.email}
                            onChange={(e) =>
                              handleInputChange('email', e.target.value)
                            }
                            className={errors.email ? 'border-red-500' : ''}
                          />
                          {errors.email && (
                            <p className="text-red-500 text-sm mt-1">
                              {errors.email}
                            </p>
                          )}
                        </div>
                        <div>
                          <Label htmlFor="phone">Phone Number *</Label>
                          <Input
                            id="phone"
                            value={formData.phone}
                            onChange={(e) =>
                              handleInputChange('phone', e.target.value)
                            }
                            className={errors.phone ? 'border-red-500' : ''}
                          />
                          {errors.phone && (
                            <p className="text-red-500 text-sm mt-1">
                              {errors.phone}
                            </p>
                          )}
                        </div>
                      </div>
                    </div>

                    {/* Rental Details */}
                    <div>
                      <h3 className="text-lg font-semibold mb-4 text-gray-900">
                        Rental Details
                      </h3>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <Label htmlFor="pickupDate">Pickup Date *</Label>
                          <Input
                            id="pickupDate"
                            type="date"
                            value={formData.pickupDate}
                            onChange={(e) =>
                              handleInputChange('pickupDate', e.target.value)
                            }
                            className={
                              errors.pickupDate ? 'border-red-500' : ''
                            }
                          />
                          {errors.pickupDate && (
                            <p className="text-red-500 text-sm mt-1">
                              {errors.pickupDate}
                            </p>
                          )}
                        </div>
                        <div>
                          <Label htmlFor="returnDate">Return Date *</Label>
                          <Input
                            id="returnDate"
                            type="date"
                            value={formData.returnDate}
                            onChange={(e) =>
                              handleInputChange('returnDate', e.target.value)
                            }
                            className={
                              errors.returnDate ? 'border-red-500' : ''
                            }
                          />
                          {errors.returnDate && (
                            <p className="text-red-500 text-sm mt-1">
                              {errors.returnDate}
                            </p>
                          )}
                        </div>
                        <div className="md:col-span-2">
                          <Label htmlFor="pickupLocation">
                            Pickup Location *
                          </Label>
                          <Input
                            id="pickupLocation"
                            placeholder="Enter pickup address or location"
                            value={formData.pickupLocation}
                            onChange={(e) =>
                              handleInputChange(
                                'pickupLocation',
                                e.target.value
                              )
                            }
                            className={
                              errors.pickupLocation ? 'border-red-500' : ''
                            }
                          />
                          {errors.pickupLocation && (
                            <p className="text-red-500 text-sm mt-1">
                              {errors.pickupLocation}
                            </p>
                          )}
                        </div>
                        <div className="md:col-span-2">
                          <Label htmlFor="driverLicense">
                            Driver License Number *
                          </Label>
                          <Input
                            id="driverLicense"
                            placeholder="Enter your driver license number"
                            value={formData.driverLicense}
                            onChange={(e) =>
                              handleInputChange('driverLicense', e.target.value)
                            }
                            className={
                              errors.driverLicense ? 'border-red-500' : ''
                            }
                          />
                          {errors.driverLicense && (
                            <p className="text-red-500 text-sm mt-1">
                              {errors.driverLicense}
                            </p>
                          )}
                        </div>
                        <div className="md:col-span-2">
                          <Label htmlFor="additionalRequests">
                            Additional Requests (Optional)
                          </Label>
                          <textarea
                            id="additionalRequests"
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            rows={3}
                            placeholder="Any special requests or requirements..."
                            value={formData.additionalRequests}
                            onChange={(e) =>
                              handleInputChange(
                                'additionalRequests',
                                e.target.value
                              )
                            }
                          />
                        </div>
                      </div>
                    </div>

                    <Button
                      type="submit"
                      size="lg"
                      className="w-full bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-semibold py-4 text-lg rounded-xl transition-all duration-300"
                    >
                      Submit Quote Request
                    </Button>
                  </form>
                </CardContent>
              </Card>
            </div>

            {/* Summary */}
            <div className="lg:col-span-1">
              <Card className="border-0 shadow-lg sticky top-4">
                <CardHeader>
                  <CardTitle className="text-xl">Rental Summary</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex items-center space-x-4">
                    <img
                      src={car.image}
                      alt={car.name}
                      className="w-16 h-12 object-cover rounded-lg"
                    />
                    <div>
                      <h4 className="font-semibold">
                        {car.brand} {car.name}
                      </h4>
                      <p className="text-sm text-gray-600">{car.type}</p>
                    </div>
                  </div>

                  <div className="border-t pt-4 space-y-2">
                    <div className="flex justify-between">
                      <span>Daily Rate:</span>
                      <span className="font-semibold">${car.pricePerDay}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Duration:</span>
                      <span className="font-semibold">
                        {days} {days === 1 ? 'day' : 'days'}
                      </span>
                    </div>
                    <div className="border-t pt-2 flex justify-between text-lg font-bold">
                      <span>Total:</span>
                      <span className="text-blue-600">${totalPrice}</span>
                    </div>
                  </div>

                  <div className="text-xs text-gray-500 space-y-1">
                    <p>* Price includes basic insurance</p>
                    <p>* Additional fees may apply</p>
                    <p>* Final price confirmed upon booking</p>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Quote;
