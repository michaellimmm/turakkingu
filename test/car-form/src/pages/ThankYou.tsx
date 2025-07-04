import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

const ThankYou = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-center">
            <div className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-blue-700 rounded-lg"></div>
              <h1 className="text-2xl font-bold text-gray-900">DriveSelect</h1>
            </div>
          </div>
        </div>
      </header>

      <div className="container mx-auto px-4 py-16">
        <div className="max-w-2xl mx-auto text-center">
          {/* Success Animation */}
          <div className="mb-8">
            <div className="mx-auto w-24 h-24 bg-gradient-to-r from-green-400 to-green-600 rounded-full flex items-center justify-center mb-6">
              <svg
                className="w-12 h-12 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={3}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
          </div>

          <Card className="border-0 shadow-2xl">
            <CardHeader className="text-center">
              <CardTitle className="text-4xl font-bold text-gray-900 mb-4">
                Thank You!
              </CardTitle>
              <p className="text-xl text-gray-600">
                Your quote request has been successfully submitted
              </p>
            </CardHeader>

            <CardContent className="space-y-6">
              <div className="bg-blue-50 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-blue-900 mb-3">
                  What happens next?
                </h3>
                <div className="text-left space-y-3">
                  <div className="flex items-center space-x-3">
                    <div className="w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center text-sm font-bold">
                      1
                    </div>
                    <p className="text-gray-700">
                      Our team will review your quote request within 2 hours
                    </p>
                  </div>
                  <div className="flex items-center space-x-3">
                    <div className="w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center text-sm font-bold">
                      2
                    </div>
                    <p className="text-gray-700">
                      You'll receive a detailed quote via email
                    </p>
                  </div>
                  <div className="flex items-center space-x-3">
                    <div className="w-8 h-8 bg-blue-600 text-white rounded-full flex items-center justify-center text-sm font-bold">
                      3
                    </div>
                    <p className="text-gray-700">
                      Confirm your booking and complete the reservation
                    </p>
                  </div>
                </div>
              </div>

              <div className="bg-green-50 rounded-lg p-4">
                <p className="text-green-800 font-medium">
                  📧 A confirmation email has been sent to your inbox with your
                  reference number.
                </p>
              </div>

              <div className="space-y-4">
                <a href="http://cardealer.local/" className="block">
                  <Button
                    size="lg"
                    className="w-full bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white font-semibold py-4 text-lg rounded-xl transition-all duration-300"
                  >
                    Browse More Cars
                  </Button>
                </a>

                <div className="text-center">
                  <p className="text-gray-600 mb-2">
                    Need immediate assistance?
                  </p>
                  <p className="text-blue-600 font-semibold">
                    Call us at: (555) 123-4567
                  </p>
                  <p className="text-blue-600 font-semibold">
                    Email: support@driveselect.com
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-12 mt-16">
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

export default ThankYou;
