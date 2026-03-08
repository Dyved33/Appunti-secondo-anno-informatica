namespace PortalStudioUnicaAzione.Tests
{
    public class UrlServiceTests
    {
        //private readonly PortalStudioUnicaAzioneContext _context;
        //private readonly IShortUrlService _urlService;

        //public UrlServiceTests()
        //{
        //    var options = new DbContextOptionsBuilder<PortalStudioUnicaAzioneContext>()
        //        .UseInMemoryDatabase(databaseName: Uuid7.NewGuid().ToString())
        //        .Options;

        //    _context = new PortalStudioUnicaAzioneContext(options);
        //    _urlService = new ShortUrlService(_context);
        //}

        //[Fact]
        //public async Task ShortenUrlAsync_ValidUrl_ReturnsUrlWithShortCode()
        //{
        //    // Arrange
        //    var originalUrl = "https://www.google.com";

        //    // Act
        //    var result = await _urlService.ShortenUrlAsync(originalUrl);

        //    // Assert
        //    Assert.NotNull(result);
        //    Assert.Equal(originalUrl, result.OriginalUrl);
        //    Assert.NotEmpty(result.ShortCode);
        //    Assert.True(result.ShortCode.Length >= 5 && result.ShortCode.Length <= 8);
        //}

        //[Fact]
        //public async Task ShortenUrlAsync_DuplicateUrl_ReturnsSameUrl()
        //{
        //    // Arrange
        //    var originalUrl = "https://www.github.com";
        //    var firstResult = await _urlService.ShortenUrlAsync(originalUrl);

        //    // Act
        //    var secondResult = await _urlService.ShortenUrlAsync(originalUrl);

        //    // Assert
        //    Assert.Equal(firstResult.Id, secondResult.Id);
        //    Assert.Equal(firstResult.ShortCode, secondResult.ShortCode);
        //}

        //[Theory]
        //[InlineData("https://www.google.com")]
        //[InlineData("http://github.com")]
        //[InlineData("www.stackoverflow.com")]
        //[InlineData("example.com")]
        //public void IsValidUrl_ValidUrls_ReturnsTrue(string url)
        //{
        //    // Act
        //    var result = _urlService.IsValidUrl(url);

        //    // Assert
        //    Assert.True(result);
        //}

        //[Theory]
        //[InlineData("")]
        //[InlineData("not-a-url")]
        //[InlineData("ftp://example.com")]
        //[InlineData("javascript:alert('test')")]
        //public void IsValidUrl_InvalidUrls_ReturnsFalse(string url)
        //{
        //    // Act
        //    var result = _urlService.IsValidUrl(url);

        //    // Assert
        //    Assert.False(result);
        //}

        //[Fact]
        //public async Task GetUrlByShortCodeAsync_ExistingCode_ReturnsUrl()
        //{
        //    // Arrange
        //    var originalUrl = "https://www.example.com";
        //    var createdUrl = await _urlService.ShortenUrlAsync(originalUrl);

        //    // Act
        //    var result = await _urlService.GetUrlByShortCodeAsync(createdUrl.ShortCode);

        //    // Assert
        //    Assert.NotNull(result);
        //    Assert.Equal(originalUrl, result.OriginalUrl);
        //    Assert.Equal(createdUrl.ShortCode, result.ShortCode);
        //}

        //[Fact]
        //public async Task GetUrlByShortCodeAsync_NonExistingCode_ReturnsNull()
        //{
        //    // Act
        //    var result = await _urlService.GetUrlByShortCodeAsync("nonexistent");

        //    // Assert
        //    Assert.Null(result);
        //}

        //[Fact]
        //public async Task IncrementClickCountAsync_ExistingUrl_IncrementsCount()
        //{
        //    // Arrange
        //    var originalUrl = "https://www.test.com";
        //    var createdUrl = await _urlService.ShortenUrlAsync(originalUrl);
        //    var initialCount = createdUrl.ClickCount;

        //    // Act
        //    await _urlService.IncrementClickCountAsync(createdUrl.ShortCode);

        //    // Assert
        //    var updatedUrl = await _urlService.GetUrlByShortCodeAsync(createdUrl.ShortCode);
        //    Assert.NotNull(updatedUrl);
        //    Assert.Equal(initialCount + 1, updatedUrl.ClickCount);
        //}

        //[Fact]
        //public async Task DeleteUrlAsync_ExistingUrl_ReturnsTrue()
        //{
        //    // Arrange
        //    var originalUrl = "https://www.delete-test.com";
        //    var createdUrl = await _urlService.ShortenUrlAsync(originalUrl);

        //    // Act
        //    var result = await _urlService.DeleteUrlAsync(createdUrl.ShortCode);

        //    // Assert
        //    Assert.True(result);
            
        //    var deletedUrl = await _urlService.GetUrlByShortCodeAsync(createdUrl.ShortCode);
        //    Assert.Null(deletedUrl);
        //}

        //[Fact]
        //public async Task DeleteUrlAsync_NonExistingUrl_ReturnsFalse()
        //{
        //    // Act
        //    var result = await _urlService.DeleteUrlAsync("nonexistent");

        //    // Assert
        //    Assert.False(result);
        //}

        //public void Dispose()
        //{
        //    _context.Dispose();
        //}
    }
}