class Product {
  final String id;
  final String name;
  final double price;
  final String currency;
  final List<String> imagePaths;
  final String sourceUrl;
  final String companyName;
  final String scrapedAt;

  Product({
    required this.id,
    required this.name,
    required this.price,
    required this.currency,
    required this.imagePaths,
    required this.sourceUrl,
    required this.companyName,
    required this.scrapedAt,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      price: (json['price'] ?? 0).toDouble(),
      currency: json['currency'] ?? 'USD',
      imagePaths: List<String>.from(json['image_paths'] ?? []),
      sourceUrl: json['source_url'] ?? '',
      companyName: json['company_name'] ?? '',
      scrapedAt: json['scraped_at'] ?? '',
    );
  }
}
