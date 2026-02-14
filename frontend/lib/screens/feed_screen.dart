import 'package:flutter/material.dart';
import '../models/product.dart';
import '../services/product_service.dart';
import 'widgets/product_card.dart';

class FeedScreen extends StatefulWidget {
  const FeedScreen({super.key});

  @override
  State<FeedScreen> createState() => _FeedScreenState();
}

class _FeedScreenState extends State<FeedScreen> {
  final List<Product> _products = [];
  String _nextCursor = '';
  bool _isLoading = false;
  bool _hasMore = true;

  @override
  void initState() {
    super.initState();
    _loadProducts();
  }

  Future<void> _loadProducts() async {
    if (_isLoading || !_hasMore) return;

    setState(() {
      _isLoading = true;
    });

    try {
      final response = await ProductService.getProducts(
        cursor: _nextCursor.isNotEmpty ? _nextCursor : null,
      );

      setState(() {
        _products.addAll(response.products);
        _nextCursor = response.nextCursor;
        _hasMore = response.nextCursor.isNotEmpty;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_products.isEmpty && _isLoading) {
      return const Scaffold(
        backgroundColor: Colors.black,
        body: Center(
          child: CircularProgressIndicator(color: Colors.white),
        ),
      );
    }

    return Scaffold(
      backgroundColor: Colors.black,
      body: PageView.builder(
        scrollDirection: Axis.vertical,
        itemCount: _products.length,
        onPageChanged: (index) {
          // Load more when 3 products from the end
          if (index >= _products.length - 3) {
            _loadProducts();
          }
        },
        itemBuilder: (context, index) {
          return ProductCard(product: _products[index]);
        },
      ),
    );
  }
}
