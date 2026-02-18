import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/product.dart';

class ProductResponse {
  final List<Product> products;
  final String nextCursor;

  ProductResponse({required this.products, required this.nextCursor});
}

class ProductService {
  static const String baseUrl = 'http://172.23.239.68:8080';

  static Future<ProductResponse> getProducts({String? cursor, int limit = 10}) async {
    var url = '$baseUrl/api/products?limit=$limit';
    if (cursor != null && cursor.isNotEmpty) {
      url += '&cursor=$cursor';
    }

    final response = await http.get(Uri.parse(url));

    if (response.statusCode != 200) {
      throw Exception('Failed to load products');
    }

    final data = json.decode(response.body);
    final products = (data['products'] as List)
        .map((json) => Product.fromJson(json))
        .toList();

    return ProductResponse(
      products: products,
      nextCursor: data['next_cursor'] ?? '',
    );
  }
}
