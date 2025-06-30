import type { Car } from '@/types/car';

export const cars: Car[] = [
  {
    id: '1',
    name: 'Model S',
    brand: 'Tesla',
    type: 'Luxury Sedan',
    image:
      'https://images.unsplash.com/photo-1617788138017-80ad40651399?w=800&h=600&fit=crop',
    pricePerDay: 120,
    year: 2023,
    transmission: 'Automatic',
    fuelType: 'Electric',
    seats: 5,
    doors: 4,
    features: ['Autopilot', 'Premium Audio', 'Heated Seats', 'WiFi Hotspot'],
    description:
      "Experience the future of driving with Tesla's flagship sedan. Advanced autopilot features and premium comfort make every journey extraordinary.",
  },
  {
    id: '2',
    name: 'Mustang GT',
    brand: 'Ford',
    type: 'Sports Car',
    image:
      'https://images.unsplash.com/photo-1584345604476-8ec5e12e42dd?w=800&h=600&fit=crop',
    pricePerDay: 95,
    year: 2023,
    transmission: 'Manual',
    fuelType: 'Gasoline',
    seats: 4,
    doors: 2,
    features: [
      'Performance Package',
      'Premium Sound',
      'Sport Suspension',
      'Track Mode',
    ],
    description:
      'Unleash your inner speed demon with the iconic Ford Mustang GT. Raw power meets classic American muscle car heritage.',
  },
  {
    id: '3',
    name: 'X5',
    brand: 'BMW',
    type: 'Luxury SUV',
    image:
      'https://images.unsplash.com/photo-1555215695-3004980ad54e?w=800&h=600&fit=crop',
    pricePerDay: 110,
    year: 2023,
    transmission: 'Automatic',
    fuelType: 'Gasoline',
    seats: 7,
    doors: 4,
    features: [
      'All-Wheel Drive',
      'Panoramic Roof',
      'Premium Interior',
      'Advanced Safety',
    ],
    description:
      'The perfect blend of luxury and versatility. BMW X5 offers premium comfort with the capability to handle any adventure.',
  },
  {
    id: '4',
    name: 'Camry Hybrid',
    brand: 'Toyota',
    type: 'Sedan',
    image:
      'https://images.unsplash.com/photo-1621007947382-bb3c3994e3fb?w=800&h=600&fit=crop',
    pricePerDay: 65,
    year: 2023,
    transmission: 'Automatic',
    fuelType: 'Hybrid',
    seats: 5,
    doors: 4,
    features: [
      'Fuel Efficient',
      'Safety Sense 2.0',
      'Wireless Charging',
      'Smart Entry',
    ],
    description:
      'Reliable, efficient, and comfortable. The Toyota Camry Hybrid is perfect for both business trips and family adventures.',
  },
  {
    id: '5',
    name: 'Wrangler',
    brand: 'Jeep',
    type: 'SUV',
    image:
      'https://images.unsplash.com/photo-1606664515524-ed2f786a0bd6?w=800&h=600&fit=crop',
    pricePerDay: 85,
    year: 2023,
    transmission: 'Manual',
    fuelType: 'Gasoline',
    seats: 5,
    doors: 4,
    features: [
      '4x4 Capability',
      'Removable Doors',
      'Off-Road Ready',
      'Rock Rails',
    ],
    description:
      'Born for adventure. The Jeep Wrangler is your gateway to exploring the great outdoors with unmatched off-road capability.',
  },
  {
    id: '6',
    name: 'Civic',
    brand: 'Honda',
    type: 'Compact',
    image:
      'https://images.unsplash.com/photo-1606664515524-ed2f786a0bd6?w=800&h=600&fit=crop',
    pricePerDay: 55,
    year: 2023,
    transmission: 'Automatic',
    fuelType: 'Gasoline',
    seats: 5,
    doors: 4,
    features: [
      'Honda Sensing',
      'Apple CarPlay',
      'Fuel Efficient',
      'Spacious Interior',
    ],
    description:
      'The Honda Civic combines reliability, efficiency, and modern technology in a compact package perfect for city driving.',
  },
];
