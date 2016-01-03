  public class Attachment
    {
        [ElasticProperty(Name = "_content")]
        public string Content { get; set; }

        [ElasticProperty(Name = "_content_type")]
        public string ContentType { get; set; }

        [ElasticProperty(Name = "_name")]
        public string Name { get; set; }
    }